package main

import (
	"context"
	"fmt"
	"funovation_23/cmd/ci_cd/runner"
	"os"
	"strings"

	"golang.org/x/exp/slices"
	"golang.org/x/mod/modfile"
)

const (
	golangPatchVersion = "0"
	buildOption        = "build"
	testOption         = "test"
	pushOption         = "push"
)

var options = []string{
	buildOption,
	testOption,
	pushOption,
}

func main() {
	mode := os.Args[1:]
	if len(mode) == 0 {
		fmt.Printf("Usage: ci [%s]\n", strings.Join(options, "|"))
		return
	}
	if !slices.Contains(options, mode[0]) {
		fmt.Printf("Usage: ci [%s]\n", strings.Join(options, "|"))
		os.Exit(1)
		return
	}

	modFile, err := os.ReadFile("go.mod")
	if err != nil {
		os.Exit(1)
		fmt.Printf("Failed: %s\n", err.Error())
		return
	}

	file, err := modfile.Parse("go.mod", modFile, nil)
	if err != nil {
		os.Exit(1)
		fmt.Printf("Failed: %s\n", err.Error())
		return
	}

	golangVersion := file.Go.Version
	golangVersion = fmt.Sprintf("%s", golangVersion)
	fmt.Printf("Running ci with go version %v\n", golangVersion)

	if mode[0] == buildOption {
		if _, err := runner.Build(context.Background(), golangVersion); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else if mode[0] == testOption {
		if err := runner.DoTests(context.Background(), golangVersion); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else if mode[0] == pushOption {
		if len(mode) < 2 {
			fmt.Println("Usage: ci push [image]")
			return
		}
		if err := runner.PushRegistry(context.Background(), mode[1], golangVersion); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
