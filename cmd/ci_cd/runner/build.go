package runner

import (
	"context"
	"fmt"
	"os"

	"dagger.io/dagger"
)

const (
	serverBinary     = "build/start_server"
	migrationsBinary = "build/run_migrations"
)

func Build(ctx context.Context, goVersion string) (*dagger.Container, error) {
	fmt.Println("Building with Dagger")

	// initialize Dagger client
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
	if err != nil {
		return nil, err
	}
	defer client.Close()

	return BuildWithClient(ctx, client, goVersion)
}

func BuildWithClient(ctx context.Context, client *dagger.Client, goVersion string) (*dagger.Container, error) {
	fmt.Println("Building with Dagger")

	golang, err := getGolangConteinerWithDependencies(ctx, client, goVersion)
	if err != nil {
		return nil, err
	}

	// add environment variables
	golang = golang.
		WithEnvVariable("CGO_ENABLED", "0").
		WithEnvVariable("GOOS", "linux").
		WithEnvVariable("GOARCH", "amd64")

	bufCommands := []command{
		{"go", "build", "-o", serverBinary, "./cmd/main.go"},
		{"go", "build", "-o", migrationsBinary, "./cmd/migrations/main.go"},
	}

	for _, command := range bufCommands {
		golang = golang.WithExec(command)
		_, err = golang.Sync(ctx)
		if err != nil {
			return nil, err
		}
	}

	return golang, nil
}
