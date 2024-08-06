package main

import (
	"testing-releaser/app/shared/constants"
	_ "testing-releaser/app/shared/infrastructure/healthcheck"
	_ "testing-releaser/app/shared/infrastructure/observability"
	"testing-releaser/app/shared/infrastructure/serverwrapper"
	_ "embed"
	"log"
	"os"

	ioc "github.com/Ignaciojeria/einar-ioc"
)

//go:embed .version
var version string

func main() {
	os.Setenv(constants.Version, version)
	if err := ioc.LoadDependencies(); err != nil {
		log.Fatal(err)
	}
	serverwrapper.Start()
}
