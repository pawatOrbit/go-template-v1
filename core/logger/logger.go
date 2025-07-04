package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync"

	"github.com/go-slog/otelslog"
	"github.com/pawatOrbit/ai-mock-data-service/go/utils/runtime"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

var (
	Log  *zap.Logger
	Slog *slog.Logger
	m    sync.Mutex
)

var Env, ServiceName, Version string

type Field = zap.Field

type LogConfig struct {
	Env             string `mapstructure:"env"`
	ServiceName     string `mapstructure:"serviceName"`
	Level           string `mapstructure:"level"`
	UseJsonEncoder  bool   `mapstructure:"useJsonEncoder"`
	StacktraceLevel string `mapstructure:"stacktraceLevel"`
	FileEnabled     bool   `mapstructure:"fileEnabled"`
	FileSize        int    `mapstructure:"fileSize"`
	FilePath        string `mapstructure:"filePath"`
	FileCompress    bool   `mapstructure:"fileCompress"`
	MaxAge          int    `mapstructure:"maxAge"`
	MaxBackups      int    `mapstructure:"maxBackups"`
}

func InitLogger(validateProfile runtime.Environment) {
	m.Lock()
	defer m.Unlock()

	// Will try to get env from os env
	Env = os.Getenv("DD_ENV")
	ServiceName = os.Getenv("DD_SERVICE")
	Version = os.Getenv("DD_VERSION")

	Slog = newZapLogger(validateProfile)
	// Slog = newSlogLogger(validateProfile)
	// Slog = slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(Slog)
	CompileCanonicalLogTemplate()
	slog.InfoContext(context.Background(), "Logger initialized")
}

var _ slog.Handler = Handler{}

type Handler struct {
	handler slog.Handler
}

func NewOtelHandler(handler slog.Handler) Handler {
	return Handler{handler: otelslog.NewHandler(handler)}
}

func NewHandler(handler slog.Handler) Handler {
	return Handler{handler: handler}
}

func (h Handler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h Handler) Handle(ctx context.Context, record slog.Record) error {
	AddDDFields(ctx, &record)
	return h.handler.Handle(ctx, record)
}

func AddDDFields(ctx context.Context, record *slog.Record) {
	spanCtx := trace.SpanContextFromContext(ctx)
	var traceID, spanID string

	if spanCtx.HasTraceID() {
		traceID = spanCtx.TraceID().String()
		record.AddAttrs(slog.String("trace_id", traceID))
	}

	if spanCtx.HasSpanID() {
		record.AddAttrs(slog.String("span_id", spanID))
	}

	record.AddAttrs(slog.Group("dd",
		slog.String("env", Env),
		slog.String("service", ServiceName),
		slog.String("trace_id", traceID),
		slog.String("span_id", spanID),
		slog.String("version", Version),
	))
}

func (h Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return Handler{h.handler.WithAttrs(attrs)}
}

func (h Handler) WithGroup(name string) slog.Handler {
	return h.handler.WithGroup(name)
}

func getLogProfile(validateProfile runtime.Environment) LogConfig {
	switch validateProfile {
	case "local":
		return LogConfig{
			Env:             "local",
			ServiceName:     ServiceName,
			Level:           "debug",
			UseJsonEncoder:  false,
			StacktraceLevel: "error",
			FileEnabled:     false,
			FilePath:        "logs/app.log",
			FileSize:        100, // megabytes
			FileCompress:    true,
			MaxAge:          30, // days
			MaxBackups:      3,  // number of log files
		}
	default:
		return LogConfig{
			Env:             Env,
			ServiceName:     ServiceName,
			Level:           "debug",
			UseJsonEncoder:  true,
			StacktraceLevel: "error",
			FileEnabled:     false,
			FilePath:        "logs/app.log",
			FileSize:        100, // megabytes
			FileCompress:    true,
			MaxAge:          30, // days
			MaxBackups:      3,  // number of log files
		}
	}
}

type Pathfinder struct {
	svc string
}

func NewPathfinder(svc string) Pathfinder {
	return Pathfinder{svc: svc}
}

func (p Pathfinder) InfoContext(ctx context.Context, msg string, fields ...any) {
	slog.InfoContext(ctx, fmt.Sprintf("[service][%s] %s", p.svc, msg), fields...)
}

func (p Pathfinder) ErrorContext(ctx context.Context, msg string, fields ...any) {
	slog.ErrorContext(ctx, fmt.Sprintf("[service][%s] %s", p.svc, msg), fields...)
}

func (p Pathfinder) DebugContext(ctx context.Context, msg string, fields ...any) {
	slog.DebugContext(ctx, fmt.Sprintf("[service][%s] %s", p.svc, msg), fields...)
}

func (p Pathfinder) NewPathfinder(svc string) Pathfinder {
	return Pathfinder{svc: p.svc + "." + svc}
}
