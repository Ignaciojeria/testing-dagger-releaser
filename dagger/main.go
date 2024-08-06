package main

import (
	"context"
	"fmt"
	"os"
	"strings"

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

	// Read the tag from .version file
	tag, err := os.ReadFile(".version")
	if err != nil {
		return fmt.Errorf("failed to read .version file: %w", err)
	}

	// Trim any extra whitespace from the tag
	tagStr := strings.TrimSpace(string(tag))

	// Get reference to the local project
	src := client.Host().Directory(".")

	// Get `goreleaser` image
	container := client.Container().From("goreleaser/goreleaser:latest")

	// Set environment
	container = container.WithEnvVariable("GITHUB_TOKEN", os.Getenv("GITHUB_ACCESS_TOKEN"))

	// Mount cloned repository into `goreleaser` image
	container = container.WithDirectory("/src", src).WithWorkdir("/src")

	// Define the application build command with the tag from .version
	container = container.WithExec([]string{"goreleaser", "--config", ".goreleaser.yml", "--release", tagStr})

	// Get reference to build output directory in container
	output := container.Directory("dist")

	// Write contents of container build/ directory to the host
	_, err = output.Export(ctx, "dist")
	if err != nil {
		return err
	}

	return nil
}
