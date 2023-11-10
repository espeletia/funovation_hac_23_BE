package runner

import (
	"context"
	"fmt"
	"os"
	"time"

	"dagger.io/dagger"
)

func getGolangConteinerWithDependencies(ctx context.Context, client *dagger.Client, goVersion string) (*dagger.Container, error) {
	src := client.Host().Directory(".")

	postgresCache := client.CacheVolume("postgres")
	postgres := client.Container().From("postgres:15.2").
		WithEnvVariable("POSTGRES_USER", "postgres").
		WithEnvVariable("CACHEBUSTER", time.Now().String()).
		WithEnvVariable("POSTGRES_PASSWORD", "postgres").
		WithEnvVariable("POSTGRES_DB", "funovation").
		WithExposedPort(5432).
		WithMountedCache("/data", postgresCache).
		WithExec(nil).AsService()

	// create a cache volume
	goCache := client.CacheVolume("go")

	golang := client.Container().
		From(fmt.Sprintf("golang:%s-alpine", goVersion)).
		WithServiceBinding("db", postgres).
		WithEnvVariable("DATABASE_URL", "postgres://postgres:postgres@db:5432/funovation?sslmode=disable").
		WithMountedCache("~/.go", goCache)

	golang = golang.WithDirectory("/src", src, dagger.ContainerWithDirectoryOpts{
		Include: []string{
			"go.*",
		},
	}).
		WithWorkdir("/src")

	for _, command := range installDependencies {
		golang = golang.WithExec(command)
		out, err := golang.Stdout(ctx)
		if err != nil {
			return nil, err
		}
		fmt.Println(out)
	}

	golang = golang.
		WithDirectory("/src", src, dagger.ContainerWithDirectoryOpts{
			Exclude: []string{
				"cmd/ci_cd/**",
			},
		})

	for _, command := range migrationCommands {
		golang = golang.WithExec(command)
		out, err := golang.Stdout(ctx)
		if err != nil {
			return nil, err
		}
		fmt.Println(out)
	}

	for _, command := range generateCommands {
		golang = golang.WithExec(command)
		out, err := golang.Stdout(ctx)
		if err != nil {
			return nil, err
		}
		fmt.Println(out)
	}

	return golang, nil
}

type (
	command []string
	envVar  = string
)

var installDependencies []command = []command{
	{"go", "mod", "download"},
	{"go", "get", "github.com/go-jet/jet/v2/cmd/jet"},
	{"go", "get", "github.com/99designs/gqlgen"},
	{"go", "install", "github.com/pressly/goose/v3/cmd/goose@latest"},
}

var generateCommands []command = []command{
	{"go", "run", "github.com/go-jet/jet/v2/cmd/jet", "-dsn=postgres://postgres:postgres@db:5432/funovation?sslmode=disable", "-path=./internal/ports/database/gen"},
	{"go", "run", "github.com/99designs/gqlgen"},
}

var migrationCommands []command = []command{
	{"ls", "migrations"},
	{"goose", "-dir", "migrations", "postgres", "postgres://postgres:postgres@db:5432/funovation?sslmode=disable", "up"},
}

func requiredGetenv(e envVar) (string, error) {
	if value := os.Getenv(e); value != "" {
		return value, nil
	}
	return "", fmt.Errorf("%v environment variable is not set", e)
}
