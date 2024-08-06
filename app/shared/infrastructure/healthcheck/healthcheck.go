package api

import (
	"testing-releaser/app/shared/configuration"
	"testing-releaser/app/shared/infrastructure/serverwrapper"

	ioc "github.com/Ignaciojeria/einar-ioc"
	"github.com/hellofresh/health-go/v5"
	"github.com/labstack/echo/v4"
)

func init() {
	ioc.Registry(healthCheck,
		serverwrapper.NewEchoWrapper,
		configuration.NewConf)
}

// To see usage examples of the library, visit: https://github.com/hellofresh/health-go
func healthCheck(e serverwrapper.EchoWrapper, c configuration.Conf) {
	h, _ := health.New(
		health.WithComponent(health.Component{
			Name:    c.PROJECT_NAME,
			Version: c.VERSION,
		}), health.WithSystemInfo())
	e.GET("/health", echo.WrapHandler(h.Handler()))
}
