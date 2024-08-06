package configuration

import (
	"log/slog"
	"os"

	ioc "github.com/Ignaciojeria/einar-ioc"
	"github.com/joho/godotenv"
)

type EnvLoader struct {
}

func init() {
	ioc.Registry(NewEnvLoader)
}
func NewEnvLoader() EnvLoader {
	if err := godotenv.Load(); err != nil {
		slog.Warn(".env not found, loading environment from system.")
	}
	return EnvLoader{}
}

func (env EnvLoader) Get(key string) string {
	return os.Getenv(key)
}
