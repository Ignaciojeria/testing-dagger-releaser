package configuration

import (
	"testing-releaser/app/shared/constants"
	"os"
	"strings"

	ioc "github.com/Ignaciojeria/einar-ioc"
)

type Conf struct {
	PORT                        string `required:"true"`
	VERSION                     string `required:"true"`
	ENVIRONMENT                 string `required:"false"`
	GEMINI_API_KEY              string `required:"false"`
	PROJECT_NAME                string `required:"false"`
	GOOGLE_PROJECT_ID           string `required:"false"`
	OTEL_EXPORTER_OTLP_ENDPOINT string `required:"false"`
	DD_SERVICE                  string `required:"false"`
	DD_ENV                      string `required:"false"`
	DD_VERSION                  string `required:"false"`
	DD_AGENT_HOST               string `required:"false"`
	COUNTRY                     string `required:"false"`
}

func init() {
	ioc.Registry(NewConf, NewEnvLoader)
}
func NewConf(env EnvLoader) (Conf, error) {
	conf := Conf{
		PORT:                        env.Get("PORT"),
		VERSION:                     env.Get(constants.Version),
		COUNTRY:                     strings.ToUpper(env.Get("COUNTRY")),
		PROJECT_NAME:                env.Get("PROJECT_NAME"),
		GEMINI_API_KEY:              env.Get("GEMINI_API_KEY"),
		GOOGLE_PROJECT_ID:           env.Get("GOOGLE_PROJECT_ID"),
		OTEL_EXPORTER_OTLP_ENDPOINT: env.Get("OTEL_EXPORTER_OTLP_ENDPOINT"),
		ENVIRONMENT:                 env.Get("ENVIRONMENT"),
	}
	setupDatadog(&conf, env)
	if conf.DD_SERVICE != "" && conf.DD_ENV != "" &&
		conf.DD_VERSION != "" && conf.DD_AGENT_HOST != "" &&
		conf.OTEL_EXPORTER_OTLP_ENDPOINT != "" {
		conf.OTEL_EXPORTER_OTLP_ENDPOINT = conf.DD_AGENT_HOST + ":4317"
	}

	if conf.PORT == "" {
		conf.PORT = "8080"
	}

	return validateConfig(conf)
}

func setupDatadog(c *Conf, env EnvLoader) {
	os.Setenv("DD_SERVICE", c.PROJECT_NAME)
	os.Setenv("DD_ENV", c.ENVIRONMENT)
	c.DD_SERVICE = c.PROJECT_NAME
	if env.Get("DD_ENV") == "" {
		os.Setenv("DD_ENV", "unknown")
	}
	c.DD_ENV = env.Get("DD_ENV")
	c.DD_AGENT_HOST = env.Get("DD_AGENT_HOST")
	os.Setenv("DD_VERSION", c.VERSION)
	c.DD_VERSION = c.VERSION
}
