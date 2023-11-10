package runner

import (
	"context"
	"fmt"
	"os"

	"dagger.io/dagger"
)

func DoTests(ctx context.Context, goVersion string) error {
	fmt.Println("Testing with Dagger")

	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
	if err != nil {
		return err
	}
	defer client.Close()

	golang, err := getGolangConteinerWithDependencies(ctx, client, goVersion)
	if err != nil {
		return err
	}

	golang = golang.WithExec([]string{"go", "test", "./..."})
	_, err = golang.Sync(ctx)
	if err != nil {
		return err
	}
	return nil
}
