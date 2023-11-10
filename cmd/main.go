package main

//go run github.com/pressly/goose/v3/cmd/goose postgres postgres://postgres:postgres@postgresql-funovation:5432/funovation?sslmode=disable up

//go:generate go run github.com/go-jet/jet/v2/cmd/jet -dsn=postgres://postgres:postgres@localhost:5432/funovation?sslmode=disable -path=../internal/ports/database/gen
//go:generate go run github.com/99designs/gqlgen generate
import (
	"fmt"
	"funovation_23/cmd/server"
	"os"
)

func main() {
	err := server.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
