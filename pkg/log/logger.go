package logger

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"go-source/pkg/utils"
	"io"
	"runtime"
	"strings"
)

type Logger struct {
	logger zerolog.Logger
}

// ------------------- Logger -------------------

func (l *Logger) StackTrace() *Logger {
	stack := getFullStack()
	newLg := l.logger.With().Str(KeyFileError, stack).Logger()
	return &Logger{newLg}
}

func (l *Logger) AddTraceInfoContextRequest(ctx context.Context) *Logger {
	newLg := l.logger.With().Interface("caller", l.GetCaller()).Logger()
	traceInfo := utils.GetRequestIdByContext(ctx)
	if traceInfo != nil {
		newLg = newLg.With().Interface(utils.KeyTraceInfo, traceInfo).Logger()
	}
	return &Logger{newLg}
}

func (l Logger) Output(w io.Writer) Logger {
	return Logger{l.logger.Output(w)}
}

func (l Logger) Level(lvl zerolog.Level) Logger {
	return Logger{l.logger.Level(lvl)}
}

func (l Logger) Sample(s zerolog.Sampler) Logger {
	return Logger{l.logger.Sample(s)}
}

func (l Logger) Hook(hooks ...zerolog.Hook) Logger {
	return Logger{l.logger.Hook(hooks...)}
}

// ------------------- Logger -------------------

// ------------------- Context -------------------

func (l Logger) With() Context {
	return Context{l: l}
}

// ------------------- Context -------------------

// ------------------- context.Context -------------------

func (l Logger) WithContext(ctx context.Context) context.Context {
	return l.logger.WithContext(ctx)
}

// ------------------- context.Context -------------------

// ------------------- Event -------------------

func (l *Logger) Trace() *Event {
	return &Event{l.logger.Trace()}
}

func (l *Logger) Debug() *Event {
	return &Event{l.logger.Debug()}
}

func (l *Logger) Info() *Event {
	return &Event{l.logger.Info()}
}

func (l *Logger) Warn() *Event {
	return &Event{l.logger.Warn()}
}

func (l *Logger) Error() *Event {
	return &Event{l.logger.Error()}
}

func (l *Logger) Err(err error) *Event {
	return &Event{l.logger.Err(err)}
}

func (l *Logger) Fatal() *Event {
	return &Event{l.logger.Fatal()}
}

func (l *Logger) Panic() *Event {
	return &Event{l.logger.Panic()}
}

func (l *Logger) WithLevel(level zerolog.Level) *Event {
	return &Event{l.logger.WithLevel(level)}
}

func (l *Logger) Log() *Event {
	return &Event{l.logger.Log()}
}

// ------------------- Event -------------------

// ------------------- Extend -------------------

func (l *Logger) GetCaller() string {
	pc, file, line, ok := runtime.Caller(2) // Adjust the call stack index as needed
	if !ok {
		return ""
	}

	fullFnName := runtime.FuncForPC(pc).Name()
	parts := strings.Split(fullFnName, ".")
	fnName := parts[len(parts)-1]

	return fmt.Sprintf("%s:%d %s", file, line, fnName)
}

func (l Logger) GetLevel() zerolog.Level {
	return l.logger.GetLevel()
}

func (l Logger) Write(p []byte) (n int, err error) {
	return l.logger.Write(p)
}

func (l *Logger) UpdateContext(update func(c zerolog.Context) zerolog.Context) {
	l.logger.UpdateContext(update)
}

func (l *Logger) Print(v ...interface{}) {
	l.logger.Print(v...)
}

func (l *Logger) Printf(format string, v ...interface{}) {
	l.logger.Printf(format, v...)
}

func (l *Logger) Println(v ...interface{}) {
	l.logger.Println(v...)
}
