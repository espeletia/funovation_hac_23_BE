package runner

import (
	"context"
	"fmt"
	"os"

	"dagger.io/dagger"
)

// REGISTRY_USERNAME='espeletia' REGISTRY_PASSWORD='dckr_pat_lJvrK9NqBiQclf35pa66D_WnsEA' SHORT_SHA='test' go run cmd/ci_cd/main.go push
const (
	serverEntry             = "./start_server"
	buildDirectory          = "./build/"
	configurationsDirectory = "./configurations/"
	migrationsDirectory     = "./migrations/"
	production              = "-prod"

	shortShaENV         envVar = "SHORT_SHA"
	registryUsernameEnv envVar = "REGISTRY_USERNAME"
	registryPasswordEnv envVar = "REGISTRY_PASSWORD"
)

func PushRegistry(ctx context.Context, enviroment string, goVersion string) error {
	fmt.Printf("Builing docker images\n")

	// initialize Dagger client
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
	if err != nil {
		return err
	}
	defer client.Close()

	// Build the go binary
	buildContainer, err := BuildWithClient(ctx, client, goVersion)
	if err != nil {
		return err
	}

	// get reference to the local project
	buildDir := buildContainer.Directory(buildDirectory)
	configurationsDir := buildContainer.Directory(configurationsDirectory)
	migrationsDir := buildContainer.Directory(migrationsDirectory)

	// Deploy image definition
	deployImage := client.Container().
		From("alpine:latest").
		WithDirectory("/app", buildDir).
		WithDirectory("/app/configurations", configurationsDir).
		WithDirectory("/app/migrations", migrationsDir).
		WithWorkdir("/app").
		WithEntrypoint([]string{serverEntry})

	// define ENV variables
	password, err := requiredGetenv(registryPasswordEnv)
	if err != nil {
		return err
	}
	passwordSecret := client.SetSecret("password", password)

	username, err := requiredGetenv(registryUsernameEnv)
	if err != nil {
		return err
	}

	sha, err := requiredGetenv(shortShaENV)
	if err != nil {
		return err
	}

	withAuth := deployImage.
		WithRegistryAuth("docker.io", username, passwordSecret)

	tags := []string{"latest", enviroment, sha}
	// push to registry
	prod := ""
	if enviroment == "production" {
		prod = production
	}
	for _, tag := range tags {
		address, err := withAuth.
			Publish(ctx, fmt.Sprintf("%s/funovation%s:%s", username, prod, tag))
		if err != nil {
			return err
		}
		fmt.Println("Image published at:", address)
	}

	return nil
}
