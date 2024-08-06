package serverwrapper

import (
	"testing-releaser/app/shared/configuration"
	"testing-releaser/app/shared/infrastructure/observability"
	"testing-releaser/app/shared/infrastructure/shutdown"
	"testing-releaser/app/shared/logging"
	"testing-releaser/app/shared/validator"
	"context"
	"fmt"
	"log"
	"log/slog"
	"time"

	ioc "github.com/Ignaciojeria/einar-ioc"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
)

type EchoWrapper struct {
	*echo.Echo
	conf configuration.Conf
}

func init() {
	ioc.Registry(echo.New)
	ioc.Registry(
		NewEchoWrapper,
		echo.New,
		configuration.NewConf,
		logging.NewLogger,
		validator.NewValidator)
}

func NewEchoWrapper(
	e *echo.Echo,
	c configuration.Conf,
	l logging.Logger,
	validator *validator.Validator) EchoWrapper {
	e.Validator = validator
	e.Use(otelecho.Middleware(c.PROJECT_NAME))
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogError:    true,
		HandleError: true, // forwards error to the global error handler, so it can decide appropriate status code
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			spanCtx, span := observability.Tracer.Start(c.Request().Context(), "RequestLogger")
			defer span.End()
			if v.Error == nil {
				l.SpanLogger(span).LogAttrs(spanCtx, slog.LevelInfo, "REQUEST",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
				)
			} else {
				l.SpanLogger(span).LogAttrs(spanCtx, slog.LevelError, "REQUEST_ERROR",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
					slog.String("err", v.Error.Error()),
				)
			}
			return nil
		},
	}))
	ctx, cancel := context.WithCancel(context.Background())
	shutdown.Handler(ctx, e.Shutdown, time.Second*5, cancel)
	return EchoWrapper{
		conf: c,
		Echo: e,
	}
}

func Start() error {
	return ioc.Get[EchoWrapper](NewEchoWrapper).start()
}

func (s EchoWrapper) start() error {
	s.printRoutes()
	err := s.Echo.Start(":" + s.conf.PORT)
	fmt.Println("waiting for resources to shut down....")
	time.Sleep(2 * time.Second)
	fmt.Println("done.")
	return err
}

func (s EchoWrapper) printRoutes() {
	routes := s.Echo.Routes()
	for _, route := range routes {
		log.Printf("Method: %s, Path: %s, Name: %s\n", route.Method, route.Path, route.Name)
	}
}
