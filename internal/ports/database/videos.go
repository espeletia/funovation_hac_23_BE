package database

import "database/sql"

type VideoStoreInterface interface{}

type VideosStore struct {
	db *sql.DB
}

func NewVideosStore(db *sql.DB) *VideosStore {
	return &VideosStore{db: db}
}
