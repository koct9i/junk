package log

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"time"

	ytlog "go.ytsaurus.tech/library/go/core/log"
)

const (
	slogTrace = slog.LevelDebug - 4
	slogFatal = slog.LevelError + 4
)

type YTSlog struct {
	l          *slog.Logger
	skipCaller int
}

func NewYTSlog(l *slog.Logger) *YTSlog {
	if l == nil {
		l = slog.Default()
	}
	return &YTSlog{l: l}
}

// Structured logging methods.
func (y *YTSlog) Trace(msg string, fields ...ytlog.Field) { y.handle(slogTrace, msg, fields...) }
func (y *YTSlog) Debug(msg string, fields ...ytlog.Field) { y.handle(slog.LevelDebug, msg, fields...) }
func (y *YTSlog) Info(msg string, fields ...ytlog.Field)  { y.handle(slog.LevelInfo, msg, fields...) }
func (y *YTSlog) Warn(msg string, fields ...ytlog.Field)  { y.handle(slog.LevelWarn, msg, fields...) }
func (y *YTSlog) Error(msg string, fields ...ytlog.Field) { y.handle(slog.LevelError, msg, fields...) }

func (y *YTSlog) Fatal(msg string, fields ...ytlog.Field) {
	y.handle(slogFatal, msg, fields...)
	os.Exit(1)
}

// Format-style logging methods.
func (y *YTSlog) Tracef(format string, args ...any) { y.Trace(fmt.Sprintf(format, args...)) }
func (y *YTSlog) Debugf(format string, args ...any) { y.Debug(fmt.Sprintf(format, args...)) }
func (y *YTSlog) Infof(format string, args ...any)  { y.Info(fmt.Sprintf(format, args...)) }
func (y *YTSlog) Warnf(format string, args ...any)  { y.Warn(fmt.Sprintf(format, args...)) }
func (y *YTSlog) Errorf(format string, args ...any) { y.Error(fmt.Sprintf(format, args...)) }
func (y *YTSlog) Fatalf(format string, args ...any) { y.Fatal(fmt.Sprintf(format, args...)) }

// Interface conversion methods required by YTsaurus log interfaces.
func (y *YTSlog) Structured() ytlog.Structured { return y }
func (y *YTSlog) Fmt() ytlog.Fmt               { return y }
func (y *YTSlog) Logger() ytlog.Logger         { return y }

func (y *YTSlog) WithName(name string) ytlog.Logger {
	return &YTSlog{l: y.l.With("logger", name), skipCaller: y.skipCaller}
}

func (y *YTSlog) With(fields ...ytlog.Field) ytlog.Logger {
	_, attrs := fieldsToAttrs(fields)
	return &YTSlog{l: slog.New(y.l.Handler().WithAttrs(attrs)), skipCaller: y.skipCaller}
}

func (y *YTSlog) AddCallerSkip(skip int) ytlog.Logger {
	return &YTSlog{l: y.l, skipCaller: y.skipCaller + skip}
}

func (y *YTSlog) handle(level slog.Level, msg string, fields ...ytlog.Field) error {
	ctx, attrs := fieldsToAttrs(fields)
	if !y.l.Enabled(ctx, level) {
		return nil
	}

	var pcs [1]uintptr
	runtime.Callers(3+y.skipCaller, pcs[:])

	rec := slog.NewRecord(time.Now(), level, msg, pcs[0])
	rec.AddAttrs(attrs...)

	return y.l.Handler().Handle(ctx, rec)
}

func fieldsToAttrs(fields []ytlog.Field) (ctx context.Context, attrs []slog.Attr) {
	attrs = make([]slog.Attr, 0, len(fields))
	for _, f := range fields {
		switch f.Type() {
		case ytlog.FieldTypeSkip:
			continue

		case ytlog.FieldTypeContext, ytlog.FieldTypeRawContext:
			if c, ok := f.Interface().(context.Context); ok && c != nil {
				ctx = c
			}
			continue

		case ytlog.FieldTypeString:
			attrs = append(attrs, slog.String(f.Key(), f.String()))

		case ytlog.FieldTypeBoolean:
			attrs = append(attrs, slog.Bool(f.Key(), f.Bool()))

		case ytlog.FieldTypeSigned:
			attrs = append(attrs, slog.Int64(f.Key(), f.Signed()))

		case ytlog.FieldTypeUnsigned:
			attrs = append(attrs, slog.Uint64(f.Key(), f.Unsigned()))

		case ytlog.FieldTypeFloat:
			attrs = append(attrs, slog.Float64(f.Key(), f.Float()))

		case ytlog.FieldTypeTime:
			attrs = append(attrs, slog.Time(f.Key(), f.Time()))

		case ytlog.FieldTypeDuration:
			attrs = append(attrs, slog.Duration(f.Key(), f.Duration()))

		case ytlog.FieldTypeError:
			attrs = append(attrs, slog.Any(f.Key(), f.Error()))

		case ytlog.FieldTypeByteString:
			if b, ok := f.Interface().([]byte); ok {
				attrs = append(attrs, slog.String(f.Key(), string(b)))
			} else {
				attrs = append(attrs, slog.Any(f.Key(), f.Any()))
			}

		case ytlog.FieldTypeLazyCall:
			if fn, ok := f.Interface().(func() (any, error)); ok {
				v, err := fn()
				if err != nil {
					attrs = append(attrs, slog.Any(f.Key()+"_error", err))
				} else {
					attrs = append(attrs, slog.Any(f.Key(), v))
				}
			}

		default:
			attrs = append(attrs, slog.Any(f.Key(), f.Any()))
		}
	}

	return ctx, attrs
}
