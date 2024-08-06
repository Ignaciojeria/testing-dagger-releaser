package main

import (
	"context"
	"fmt"

	"testing-releaser/dagger/steps"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}
	ctx := context.Background()
	if err := steps.PublishRelease(ctx); err != nil {
		panic(err)
	}
}
