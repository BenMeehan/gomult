package sandbox

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type nsjailExecutor struct {
	path string
}

func NewNsjailExecutor(path string) *nsjailExecutor {
	return &nsjailExecutor{path: path}
}

func (e *nsjailExecutor) Name() string { return "nsjail" }

func (e *nsjailExecutor) Available() bool {
	_, err := exec.LookPath(e.path)
	return err == nil
}

func (e *nsjailExecutor) Execute(ctx context.Context, workDir string, cmd []string, stdin string, limits Limits) (*ExecResult, error) {
	if len(cmd) == 0 {
		return nil, fmt.Errorf("empty command")
	}

	args := []string{
		"--mode", "execve",
		"--time_limit", strconv.Itoa(int(limits.TimeLimit.Seconds())),
		"--rlimit_as", strconv.FormatInt(limits.MaxMemory, 10),
		"--rlimit_fsize", strconv.FormatInt(limits.MaxFileSize, 10),
		"--rlimit_nproc", strconv.Itoa(limits.MaxProcesses),
		"--iface_no_lo",
		"--user", strconv.Itoa(limits.User),
		"--group", strconv.Itoa(limits.Group),
		"--quiet",
		"--cgroup_mem_max", "0",
		"--",
	}
	args = append(args, cmd...)

	deadline, ok := ctx.Deadline()
	if !ok {
		deadline = time.Now().Add(limits.TimeLimit + 5*time.Second)
	}

	execCtx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	c := exec.CommandContext(execCtx, e.path, args...)
	c.Dir = workDir
	c.Stdin = strings.NewReader(stdin)

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
