package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"dagger.io/dagger"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
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
	tag, err := ioutil.ReadFile(".version")
	if err != nil {
		return fmt.Errorf("failed to read .version file: %w", err)
	}

	// Trim any extra whitespace from the tag
	tagStr := strings.TrimSpace(string(tag))

	// Create and push the tag using go-git
	if err := createAndPushTag(tagStr); err != nil {
		return err
	}

	// Get reference to the local project
	src := client.Host().Directory(".")

	// Get `goreleaser` image
	container := client.Container().From("goreleaser/goreleaser:latest")

	// Set environment
	container = container.WithEnvVariable("GITHUB_TOKEN", os.Getenv("GITHUB_ACCESS_TOKEN"))

	// Mount cloned repository into `goreleaser` image
	container = container.WithDirectory("/src", src).WithWorkdir("/src")

	// Define the application build command with the tag from .version
	container = container.WithExec([]string{"goreleaser", "--rm-dist", "--config", ".goreleaser.yml"})

	// Get reference to build output directory in container
	output := container.Directory("dist")

	// Write contents of container build/ directory to the host
	_, err = output.Export(ctx, "dist")
	if err != nil {
		return err
	}

	return nil
}

func createAndPushTag(tag string) error {
	// Open the existing repository
	repo, err := git.PlainOpen(".")
	if err != nil {
		return fmt.Errorf("failed to open git repository: %w", err)
	}

	// Create the tag
	_, err = repo.CreateTag(tag, plumbing.NewHash("HEAD"), nil)
	if err != nil {
		return fmt.Errorf("failed to create git tag: %w", err)
	}

	// Push the tag
	err = repo.Push(&git.PushOptions{
		Auth: &http.TokenAuth{
			Token: os.Getenv("GITHUB_ACCESS_TOKEN"),
		},
		RefSpecs: []config.RefSpec{"+refs/tags/*:refs/tags/*"},
	})
	if err != nil {
		return fmt.Errorf("failed to push git tag: %w", err)
	}

	return nil
}
