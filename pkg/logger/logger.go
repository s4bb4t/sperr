package logger

import (
	"fmt"
	"github.com/s4bb4t/lighthouse/internal/hooks"
	"github.com/s4bb4t/lighthouse/pkg/core/levels"
	"github.com/s4bb4t/lighthouse/pkg/core/sperror"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
)

const (
	Local = "local"
	Dev   = "dev"
	Prod  = "prod"
)

type Logger struct {
	pd    func(layers int) string
	log   *slog.Logger
	stage string
	lg    string
	noop  bool
}

// Noop - creates new Logger that does nothing
//
// # It's useful for testing purposes
//
// Example:
//
//	l := logger.Noop()
//	l.Error(errors.New("test error"), logger.LevelHighDebug)
func Noop() *Logger {
	l := &Logger{noop: true}
	return l
}

// New - creates new Logger
//
// stage - one of Local, Dev, Prod
// lg - language code
// out - io.Writer to write logs to
//
// Logger's language is used only to define sperror.Error's message
func New(stage, lg string, out io.Writer) *Logger {
	if out == nil {
		out = os.Stdout
	}

	l := &Logger{lg: lg, stage: stage, pd: func(layers int) string {
		_, file, line, ok := runtime.Caller(layers + 1)
		if ok {
			absPath, err := filepath.Abs(file)
			if err != nil {
				panic(err)
			}
			return fmt.Sprintf("%s:%d", absPath, line)
		}
		return "unknown"
	}}

	switch stage {
	default:
		l.log = slog.New(newPrettyHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case Dev:
		l.log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case Prod:
		l.log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	}
	return l
}

// With - adds fields to logger
//
// It's a shortcut for slog.With()
/*
func (l *Logger) With(args ...any) {
	if l.noop {
		return
	}
	l.log = l.log.With(args...)
}
*/

// Todo: add colours

//args = append(args, slog.String("log_at", l.pd(1)))
// TODO: log_at

// Warn - logs message with source path and optional error
//
// Error can be nil - it's ok
func (l *Logger) Warn(msg string, e error, args ...any) {
	if l.noop {
		return
	}
	if e != nil {
		args = append(args, hooks.Slog(sperror.Ensure(e), levels.LevelError)...)
	}
	l.log.Warn(msg, args...)
}

// ErrorWithLevel - logs error with specified level
func (l *Logger) ErrorWithLevel(e error, lvl levels.Level) {
	if l.noop || e == nil {
		return
	}
	err := sperror.Ensure(e)
	// spin-prepare and log error
	args := hooks.Slog(err, lvl)
	l.log.Error(err.Msg(l.lg), args...)
}

// Error - logs error with default Error level
//
// lvl - error level
func (l *Logger) Error(e error) {
	if l.noop || e == nil {
		return
	}
	err := sperror.Ensure(e)
	// spin-prepare and log error
	args := hooks.Slog(err, levels.LevelError)
	l.log.Error(err.Msg(l.lg), args...)
}

// Debug - prints additional debug log to Logger's out
func (l *Logger) Debug(msg string, args ...any) {
	if l.noop {
		return
	}
	l.log.Debug(msg, args...)
}

// Info - prints additional info to Logger's out
func (l *Logger) Info(msg string, args ...any) {
	if l.noop {
		return
	}
	l.log.Info(msg, args...)
}
