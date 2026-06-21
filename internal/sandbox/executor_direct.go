package sandbox

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"
)

type directExecutor struct{}

func NewDirectExecutor() *directExecutor {
	return &directExecutor{}
}

func (e *directExecutor) Name() string     { return "direct" }
func (e *directExecutor) Available() bool  { return true }

func (e *directExecutor) Execute(ctx context.Context, workDir string, cmd []string, stdin string, limits Limits) (*ExecResult, error) {
	if len(cmd) == 0 {
		return nil, fmt.Errorf("empty command")
	}

	deadline, ok := ctx.Deadline()
	if !ok {
		deadline = time.Now().Add(limits.TimeLimit + 5*time.Second)
	}

	execCtx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	c := exec.CommandContext(execCtx, cmd[0], cmd[1:]...)
	c.Dir = workDir
	c.Stdin = bytes.NewBufferString(stdin)

	if os.Getuid() == 0 && limits.User != 0 {
		c.SysProcAttr = &syscall.SysProcAttr{
			Credential: &syscall.Credential{
				Uid: uint32(limits.User),
				Gid: uint32(limits.Group),
			},
		}
	}

	var stdout, stderr bytes.Buffer
	c.Stdout = &stdout
	c.Stderr = &stderr

	start := time.Now()
	err := c.Run()
	duration := time.Since(start)

	result := &ExecResult{
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
		Duration: duration,
	}

	if err != nil {
		if execCtx.Err() == context.DeadlineExceeded {
			result.TimedOut = true
			result.ExitCode = -1
			return result, nil
		}
		if exitErr, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitErr.ExitCode()
		} else {
			result.ExitCode = -1
		}
		return result, err
	}

	result.ExitCode = 0
	return result, nil
}
