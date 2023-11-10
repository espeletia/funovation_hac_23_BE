package database

import (
	"context"
	"database/sql"
	"funovation_23/internal/domain"
	"funovation_23/internal/ports/database/gen/funovation/public/model"
	"funovation_23/internal/ports/database/gen/funovation/public/table"

	"github.com/go-jet/jet/v2/postgres"
)

type VideoStoreInterface interface {
	CreateVideo(ctx context.Context, video *domain.DownloadedYTVideo, path string) (*domain.YoutubeVideo, error)
	GetVideo(ctx context.Context, id int64) (*domain.YoutubeVideo, error)
	GetAllVideos(ctx context.Context) ([]*domain.YoutubeVideo, error)
}

type VideosStore struct {
	db *sql.DB
}

func NewVideosStore(db *sql.DB) *VideosStore {
	return &VideosStore{db: db}
}

func (vs *VideosStore) CreateVideo(ctx context.Context, video *domain.DownloadedYTVideo, path string) (*domain.YoutubeVideo, error) {
	inssertModel := model.Videos{
		YoutubeID: video.YoutubeID,
		Title:     video.Title,
		URL:       path,
	}
	stmt := table.Videos.INSERT(
		table.Videos.YoutubeID,
		table.Videos.Title,
		table.Videos.URL,
	).
		MODEL(inssertModel).
		RETURNING(table.Videos.AllColumns)

	dest := model.Videos{}
	err := stmt.QueryContext(ctx, vs.db, &dest)
	if err != nil {
		return nil, err
	}
	return mapDBVideo(dest), nil
}

func (vs *VideosStore) GetVideo(ctx context.Context, id int64) (*domain.YoutubeVideo, error) {
	stmt := table.Videos.SELECT(
		table.Videos.AllColumns,
	).
		WHERE(table.Videos.ID.EQ(postgres.Int(id)))

	dest := []model.Videos{}
	err := stmt.QueryContext(ctx, vs.db, &dest)
	if err != nil {
		return nil, err
	}
	if len(dest) < 1 {
		return nil, nil
	}
	return mapDBVideo(dest[0]), nil
}

func (vs *VideosStore) GetAllVideos(ctx context.Context) ([]*domain.YoutubeVideo, error) {
	stmt := table.Videos.SELECT(
		table.Videos.AllColumns,
	)

	dest := []model.Videos{}
	err := stmt.QueryContext(ctx, vs.db, &dest)
	if err != nil {
		return nil, err
	}
	if len(dest) < 1 {
		return nil, nil
	}
	var result []*domain.YoutubeVideo
	for _, video := range dest {
		result = append(result, mapDBVideo(video))
	}
	return result, nil
}

func mapDBVideo(video model.Videos) *domain.YoutubeVideo {
	return &domain.YoutubeVideo{
		ID:        int64(video.ID),
		YoutubeID: video.YoutubeID,
		Title:     video.Title,
		S3Path:    video.URL,
		Status:    int64(video.Status),
	}
}
