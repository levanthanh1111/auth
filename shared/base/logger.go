package base

import (
	"github.com/rs/zerolog"
	"github.com/tpp/msf/shared/context"
)

type Logger interface {
	Trace(ctx context.Context) *zerolog.Event
	Debug(ctx context.Context) *zerolog.Event
	Info(ctx context.Context) *zerolog.Event
	Warn(ctx context.Context) *zerolog.Event
	Error(ctx context.Context) *zerolog.Event
	Faltal(ctx context.Context) *zerolog.Event
	Panic(ctx context.Context) *zerolog.Event
}

type logger struct {
	*zerolog.Logger
}

func (l *logger) Trace(ctx context.Context) *zerolog.Event {
	return l.Logger.Info().Str("req_id", ctx.ReqID())
}

func (l *logger) Debug(ctx context.Context) *zerolog.Event {
	return l.Logger.Info().Str("req_id", ctx.ReqID())
}

func (l *logger) Info(ctx context.Context) *zerolog.Event {
	return l.Logger.Info().Str("req_id", ctx.ReqID())
}

func (l *logger) Warn(ctx context.Context) *zerolog.Event {
	return l.Logger.Warn().Str("req_id", ctx.ReqID())
}

func (l *logger) Error(ctx context.Context) *zerolog.Event {
	return l.Logger.Error().Str("req_id", ctx.ReqID())
}

func (l *logger) Faltal(ctx context.Context) *zerolog.Event {
	return l.Logger.Fatal().Str("req_id", ctx.ReqID())
}

func (l *logger) Panic(ctx context.Context) *zerolog.Event {
	return l.Logger.Panic().Str("req_id", ctx.ReqID())
}

func newBaseLogger(l zerolog.Logger) Logger {
	return &logger{&l}
}
