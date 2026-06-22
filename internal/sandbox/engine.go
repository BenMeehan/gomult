package sandbox

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/benmeehan/gomult/internal/config"
)

type Engine struct {
	executor  Executor
	languages map[string]config.LangConfig
	limits    Limits
	maxSize   int64
}

func NewEngine(cfg *config.Config) (*Engine, error) {
	limits := Limits{
		TimeLimit:    cfg.Sandbox.TimeLimitDuration(),
		MaxMemory:    cfg.Sandbox.MaxMemory,
		MaxFileSize:  cfg.Sandbox.MaxFileSize,
		MaxProcesses: cfg.Sandbox.MaxProcesses,
		User:         cfg.Sandbox.RuntimeUser,
		Group:        cfg.Sandbox.RuntimeGroup,
	}

	var executor Executor
	nsjail := NewNsjailExecutor(cfg.Sandbox.NsjailPath)
	direct := NewDirectExecutor()

	switch cfg.Sandbox.Mode {
	case "nsjail":
		if !nsjail.Available() {
			return nil, fmt.Errorf("nsjail not found at %s", cfg.Sandbox.NsjailPath)
		}
		executor = nsjail
		log.Printf("sandbox: using nsjail executor")
	case "direct":
		executor = direct
		log.Printf("sandbox: using direct executor (no sandbox)")
	case "auto":
		if nsjail.Available() {
			executor = nsjail
			log.Printf("sandbox: using nsjail executor (auto-detected)")
		} else {
			executor = direct
			log.Printf("sandbox: nsjail not found, using direct executor")
		}
	default:
		return nil, fmt.Errorf("unknown sandbox mode: %s", cfg.Sandbox.Mode)
	}

	if len(cfg.Languages) == 0 {
		return nil, fmt.Errorf("no languages configured")
	}

	return &Engine{
		executor:  executor,
		languages: cfg.Languages,
		limits:    limits,
		maxSize:   cfg.Server.MaxCodeSize,
	}, nil
}

func (e *Engine) Languages() map[string]string {
	result := make(map[string]string, len(e.languages))
	for k, v := range e.languages {
		result[k] = v.Name
	}
	return result
}

func (e *Engine) Execute(langKey, code, stdin string) *ExecuteResult {
	lang, ok := e.languages[langKey]
	if !ok {
		keys := e.supportedKeys()
		return &ExecuteResult{
			Status: StatusUnsupportedLanguage,
			Output: fmt.Sprintf("unsupported language: %s (supported: %s)", langKey, keys),
		}
	}

	if int64(len(code)) > e.maxSize {
		return &ExecuteResult{
			Status: StatusCodeTooLarge,
			Output: fmt.Sprintf("code too large (max %d bytes)", e.maxSize),
		}
	}

	workDir, err := os.MkdirTemp("", "gomult-*")
	if err != nil {
		return &ExecuteResult{
			Status: StatusInternalError,
			Output: fmt.Sprintf("failed to create workdir: %v", err),
		}
	}
	defer os.RemoveAll(workDir)

	srcName := "code" + lang.Extension
	if lang.SourceFile != "" {
		srcName = lang.SourceFile
	}

	srcPath := filepath.Join(workDir, srcName)
	if err := os.WriteFile(srcPath, []byte(code), 0644); err != nil {
		return &ExecuteResult{
			Status: StatusInternalError,
			Output: fmt.Sprintf("failed to write code: %v", err),
		}
	}

	var runCmd []string
	if len(lang.CompileCmd) > 0 {
		outputPath := filepath.Join(workDir, "output")
		compileCmd := resolveCmd(lang.CompileCmd, srcPath, outputPath, workDir)

		compileCtx, cancel := context.WithTimeout(context.Background(), e.limits.TimeLimit+10*time.Second)
		defer cancel()

		compileResult, err := runCommand(compileCtx, workDir, compileCmd, "")
		if err != nil || compileResult.ExitCode != 0 {
			msg := compileResult.Stderr
			if msg == "" {
				msg = compileResult.Stdout
			}
			if msg == "" && err != nil {
				msg = err.Error()
			}
			return &ExecuteResult{
				Status: StatusCompileError,
				Output: msg,
			}
		}

		runCmd = resolveCmd(lang.RunCmd, srcPath, outputPath, workDir)
	} else {
		runCmd = resolveCmd(lang.RunCmd, srcPath, "", workDir)
	}

	ctx, cancel := context.WithTimeout(context.Background(), e.limits.TimeLimit+5*time.Second)
	defer cancel()

	result, err := e.executor.Execute(ctx, workDir, runCmd, stdin, e.limits)
	if err != nil && result.ExitCode == -1 && !result.TimedOut {
		output := result.Stderr
		if output == "" {
			output = result.Stdout
		}
		if output == "" {
			output = err.Error()
		}
		return &ExecuteResult{
			Status: StatusInternalError,
			Output: fmt.Sprintf("execution error: %v", output),
		}
	}

	if result.TimedOut {
		return &ExecuteResult{
			Status: StatusTimeLimitExceeded,
			Output: result.Stderr,
		}
	}

	if result.ExitCode != 0 {
		output := result.Stderr
		if output == "" {
			output = result.Stdout
		}
		return &ExecuteResult{
			Status:   StatusRuntimeError,
			Output:   output,
			ExitCode: result.ExitCode,
		}
	}

	return &ExecuteResult{
		Status: StatusOK,
		Output: result.Stdout,
	}
}

func (e *Engine) supportedKeys() string {
	keys := make([]string, 0, len(e.languages))
	for k := range e.languages {
		keys = append(keys, k)
	}
	return strings.Join(keys, ", ")
}

func resolveCmd(template []string, srcPath, outputPath, workDir string) []string {
	result := make([]string, len(template))
	for i, part := range template {
		s := part
		s = strings.ReplaceAll(s, "{file}", srcPath)
		s = strings.ReplaceAll(s, "{output}", outputPath)
		s = strings.ReplaceAll(s, "{dir}", workDir)
		result[i] = s
	}
	return result
}

func runCommand(ctx context.Context, workDir string, cmd []string, stdin string) (*ExecResult, error) {
	c := exec.CommandContext(ctx, cmd[0], cmd[1:]...)
	c.Dir = workDir
	if stdin != "" {
		c.Stdin = strings.NewReader(stdin)
	}

	var stdout, stderr bytes.Buffer
	c.Stdout = &stdout
	c.Stderr = &stderr

	err := c.Run()
	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			exitCode = -1
		}
	}

	return &ExecResult{
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
		ExitCode: exitCode,
	}, err
}
