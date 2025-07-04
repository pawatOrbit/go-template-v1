package logger

import (
	"context"
	"encoding/json"
	"errors"
	"html/template"
	"log/slog"
	"strings"
	"time"

	"github.com/pawatOrbit/ai-mock-data-service/go/core/exception"
)

var canonicalLogTemplate *template.Template

type Level int

const (
	Debug Level = 1 << iota
	Info
	Warn
	Error
)

type CanonicalLog struct {
	Transport string
	Traffic   string
	Method    string
	Status    int
	Path      string
	Duration  time.Duration
	Message   string
	Level     slog.Level
}

func CompileCanonicalLogTemplate() {
	logTemplate := "[{{.Transport}}][{{.Traffic}}] {{.Method}} {{.Status}} {{.Path}} {{.Duration}} - {{.Message}}"
	compiled, err := template.New("log_template").Parse(logTemplate)
	if err != nil {
		panic(err)
	}
	canonicalLogTemplate = compiled
}

func GetCanonicalLogTemplate() (*template.Template, error) {
	if canonicalLogTemplate != nil {
		return canonicalLogTemplate, nil
	}
	return nil, errors.New("canonicalLogTemplate is nil")
}

func CanonicalLogger(ctx context.Context, slogger slog.Logger, level Level, request []byte, response []byte, err error, cannonicalLog CanonicalLog, metadata []any) {
	// log the cannonical log

	logKey := cannonicalLog.Path
	var reqfields []any
	// append request log
	var jsonObj map[string]interface{}

	if unmarshalErr := json.Unmarshal(request, &jsonObj); unmarshalErr != nil {
		reqfields = append(reqfields, slog.String("request", string(request)))
	} else {
		reqfields = append(reqfields, slog.Any("request", jsonObj))
	}

	shouldSanitize := Sanitize(logKey)

	// bypass sanitization for certain paths
	shouldSanitize = false
	if shouldSanitize {
		reqfields = []any{slog.String("request", "REDACTED")}
	}

	var respFields []any
	// append response log
	if err != nil {
		level = Error
		cErr, ok := err.(*exception.ExceptionError)
		if ok && cErr != nil {
			if cErr.StackErrors != nil {
				stackTrace := exception.GetStackField(cErr.StackErrors)
				stackTraceParts := strings.Split(stackTrace.Stack, "\n\t")
				if len(stackTraceParts) > 6 {
					stackTrace.Stack = strings.Join(stackTraceParts[:6], "\n\t")
				}
				respFields = append(respFields, slog.Group("error",
					slog.String("kind", stackTrace.Kind),
					slog.String("message", stackTrace.Message),
					slog.String("stack", stackTrace.Stack),
				))
				// slogger.ErrorContext(ctx, "logger debug", slog.Bool("override", cErr.OverrideLogLevel), slog.Any("level", cErr.Level))
				if cErr.OverrideLogLevel {
					switch cErr.Level {
					case exception.LevelDebug:
						level = Debug
					case exception.LevelInfo:
						level = Info
					case exception.LevelWarn:
						level = Warn
					case exception.LevelError:
						level = Error
					default:
						level = Error
					}
				}
			}
			respFields = append(respFields, slog.Group("response",
				slog.Int("status_code", cErr.APIStatusCode),
				slog.Any("data", nil),
				slog.Group("error",
					slog.Int("code", int(cErr.Code)),
					slog.String("message", cErr.GlobalMessage),
					slog.String("debug_message", cErr.DebugMessage),
					slog.Any("details", cErr.ErrFields),
				)))
			cannonicalLog.Message = cErr.DebugMessage
		} else {
			// This is the case when the error is not an instance of ExceptionError
			var jsonObj map[string]interface{}
			if err := json.Unmarshal(response, &jsonObj); err != nil {
				respFields = append(respFields, slog.String("response",
					string(response),
				))
			} else {
				respFields = append(respFields, slog.Any("response",
					jsonObj,
				))
			}
		}
	} else {
		level = Info
		var jsonObj map[string]interface{}
		if err := json.Unmarshal(response, &jsonObj); err != nil {
			respFields = append(respFields, slog.String("response",
				string(response),
			))
		} else {
			respFields = append(respFields, slog.Any("response",
				jsonObj,
			))
		}
	}
	if shouldSanitize {
		respFields = []any{slog.String("response", "REDACTED")}
	}

	var mdFields []any
	// append md log
	mdFields = append(mdFields,
		slog.String("logger_name", "canonical"),
		slog.Group("md", metadata...),
	)

	var logMsgBuilder strings.Builder
	var logMsg string
	logTmpl, logTmplErr := GetCanonicalLogTemplate()
	if logTmplErr != nil {
		logMsg = "failed to get cannonical log template"
	} else {
		executeErr := logTmpl.Execute(&logMsgBuilder, cannonicalLog)
		if executeErr != nil {
			logMsg = "failed to execute cannonical log template"
		} else {
			logMsg = logMsgBuilder.String()
		}
	}

	fields := append(reqfields, respFields...)
	fields = append(fields, mdFields...)

	switch level {
	case Debug:
		slogger.DebugContext(ctx, logMsg, fields...)
	case Info:
		slogger.InfoContext(ctx, logMsg, fields...)
	case Warn:
		slogger.WarnContext(ctx, logMsg, fields...)
	case Error:
		slogger.ErrorContext(ctx, logMsg, fields...)
	default:
		slogger.ErrorContext(ctx, logMsg, fields...)
	}
}

var DenyPatterns = []string{
	"login",
	"api-key",
	"register",
}

func Sanitize(logKey string) bool {
	// check if deny pattern is in path
	for _, denyPattern := range DenyPatterns {
		if strings.Contains(logKey, denyPattern) {
			return true
		}
	}
	return false
}
