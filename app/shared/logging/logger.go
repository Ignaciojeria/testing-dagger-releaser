package logging

import (
	"testing-releaser/app/shared/configuration"
	"log/slog"
	"os"
	"strconv"

	ioc "github.com/Ignaciojeria/einar-ioc"
	"go.opentelemetry.io/otel/trace"
)

// Datadog trace and log correlation :
// https://docs.datadoghq.com/tracing/other_telemetry/connect_logs_and_traces/opentelemetry/?tab=go
const (
	ddTraceIDKey = "dd.trace_id"
	ddSpanIDKey  = "dd.span_id"
	ddServiceKey = "dd.service"
	ddEnvKey     = "dd.env"
	ddVersionKey = "dd.version"
)

// Default opentelemetry trace and log correlation :
const (
	traceIDKey = "trace_id"
	spanIDKey  = "span_id"
)

type Logger struct {
	*slog.Logger
	conf configuration.Conf
}

func init() {
	ioc.Registry(NewLogger, configuration.NewConf)
}
func NewLogger(conf configuration.Conf) Logger {
	return Logger{
		Logger: slog.New(slog.NewJSONHandler(os.Stdout, nil)),
		conf:   conf,
	}
}

func (l Logger) SpanLogger(span trace.Span) *slog.Logger {
	traceID := span.SpanContext().TraceID().String()
	spanID := span.SpanContext().SpanID().String()

	ddService := l.conf.DD_SERVICE
	ddEnv := l.conf.DD_ENV
	ddVersion := l.conf.DD_VERSION

	if ddService == "" || ddEnv == "" || ddVersion == "" {
		return l.Logger.With(
			slog.String(traceIDKey, traceID),
			slog.String(spanIDKey, spanID),
		)
	}
	return l.Logger.With(
		slog.String(traceIDKey, traceID),
		slog.String(spanIDKey, spanID),
		slog.String(ddTraceIDKey, convertTraceID(traceID)),
		slog.String(ddSpanIDKey, convertTraceID(spanID)),
		slog.String(ddServiceKey, ddService),
		slog.String(ddEnvKey, ddEnv),
		slog.String(ddVersionKey, ddVersion),
	)
}

func convertTraceID(id string) string {
	if len(id) < 16 {
		return ""
	}
	if len(id) > 16 {
		id = id[16:]
	}
	intValue, err := strconv.ParseUint(id, 16, 64)
	if err != nil {
		return ""
	}
	return strconv.FormatUint(intValue, 10)
}
