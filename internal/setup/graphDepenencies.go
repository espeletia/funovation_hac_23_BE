package setup

import (
	"database/sql"
	"funovation_23/graph"
	"funovation_23/internal/config"
)

func NewResolver(dbConn *sql.DB, config config.Config) (*graph.Resolver, error) {

	return &graph.Resolver{}, nil

}
