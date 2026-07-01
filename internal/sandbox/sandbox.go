package sandbox

import (
	"context"
	"time"
)

type Status int

const (
	StatusOK Status = iota
	StatusCompileError
	StatusTimeLimitExceeded
	StatusRuntimeError
	StatusMemoryLimitExceeded
	StatusInternalError
	StatusUnsupportedLanguage
	StatusCodeTooLarge
)

func (s Status) String() string {
	switch s {
	case StatusOK:
		return "ok"
	case StatusCompileError:
		return "compile_error"
	case StatusTimeLimitExceeded:
		return "timeout"
	case StatusRuntimeError:
		return "runtime_error"
	case StatusMemoryLimitExceeded:
		return "memory_limit_exceeded"
	case StatusInternalError:
		return "internal_error"
	case StatusUnsupportedLanguage:
		return "unsupported_language"
	case StatusCodeTooLarge:
		return "code_too_large"
	default:
		return "unknown"
	}
}

type Limits struct {
	TimeLimit    time.Duration
	MaxMemory    int64
	MaxFileSize  int64
	MaxProcesses int
	User         int
	Group        int
}

type ExecResult struct {
	Stdout   string
	Stderr   string
	ExitCode int
	Duration time.Duration
	TimedOut bool
}

type Executor interface {
	Execute(ctx context.Context, workDir string, cmd []string, stdin string, limits Limits) (*ExecResult, error)
	Available() bool
	Name() string
}

type ExecuteResult struct {
	Output   string
	Status   Status
	ExitCode int
}

type LangInfo struct {
	Name      string `json:"name"`
	Extension string `json:"extension"`
}
