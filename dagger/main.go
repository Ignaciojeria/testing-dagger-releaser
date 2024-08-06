package main

import (
	"context"
	"fmt"
	"os"

	"dagger.io/dagger"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}
	ctx := context.Background()
	if err := PublishRelease(ctx); err != nil {
		panic(err)
	}
}

func PublishRelease(ctx context.Context) error {
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
	defer client.Close()
	if err != nil {
		return err
	}
	// Get reference to the local project
	src := client.Host().Directory(".")

	// Get `goreleaser` image
	container := client.Container().From("goreleaser/goreleaser:latest")

	// Set environment variables
	container = container.
		WithEnvVariable("GITHUB_TOKEN", os.Getenv("GITHUB_ACCESS_TOKEN")).
		WithEnvVariable("GORELEASER_CURRENT_TAG", "1.0.6")

	// Mount cloned repository into `goreleaser` image
	container = container.WithDirectory("/src", src).WithWorkdir("/src")

	// Define the application build command
	container = container.WithExec([]string{"goreleaser", "--config", ".goreleaser.yml"})

	// Get reference to build output directory in container
	output := container.Directory("dist")

	// Write contents of container build/ directory to the host
	_, err = output.Export(ctx, "dist")
	if err != nil {
		return err
	}

	return nil
}
