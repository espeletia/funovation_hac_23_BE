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
	CreateVideo(ctx context.Context, video *domain.DownloadedYTVideo, videoPath, thumbnailPath string) (*domain.YoutubeVideo, error)
	GetVideo(ctx context.Context, id int64) (*domain.YoutubeVideo, error)
	GetAllVideos(ctx context.Context) ([]*domain.YoutubeVideo, error)
	CreateClip(ctx context.Context, clip domain.CreateClip) (*domain.Clip, error)
}

type VideosStore struct {
	db *sql.DB
}

func NewVideosStore(db *sql.DB) *VideosStore {
	return &VideosStore{db: db}
}

func (vs *VideosStore) CreateVideo(ctx context.Context, video *domain.DownloadedYTVideo, videoPath, thumbnailPath string) (*domain.YoutubeVideo, error) {
	inssertModel := model.Videos{
		YoutubeID:   video.YoutubeID,
		Title:       video.Title,
		URL:         videoPath,
		Thumbnail:   thumbnailPath,
		CustomTitle: video.CustomTitle,
		Description: video.Description,
	}
	stmt := table.Videos.INSERT(
		table.Videos.YoutubeID,
		table.Videos.Title,
		table.Videos.URL,
		table.Videos.CustomTitle,
		table.Videos.Thumbnail,
		table.Videos.Description,
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

func (vs *VideosStore) CreateClip(ctx context.Context, clip domain.CreateClip) (*domain.Clip, error) {
	insertModel := model.Clips{
		Videoid:   int32(clip.VideoID),
		URL:       clip.URL,
		Thumbnail: clip.Thumbnail,
		Order:     int32(clip.Order),
	}
	stmt := table.Clips.INSERT(
		table.Clips.AllColumns.Except(table.Clips.ID, table.Clips.CreatedAt),
	).
		MODEL(insertModel).
		RETURNING(table.Clips.AllColumns)

	dest := model.Clips{}
	err := stmt.QueryContext(ctx, vs.db, &dest)
	if err != nil {
		return nil, err
	}
	return mapDBClip(dest), nil
}

func mapDBVideo(video model.Videos) *domain.YoutubeVideo {
	return &domain.YoutubeVideo{
		ID:          int64(video.ID),
		YoutubeID:   video.YoutubeID,
		Title:       video.Title,
		S3Path:      video.URL,
		Status:      int64(video.Status),
		Thumbnail:   video.Thumbnail,
		CustomTitle: video.CustomTitle,
		Description: video.Description,
	}
}

func mapDBClip(clip model.Clips) *domain.Clip {
	return &domain.Clip{
		ID:        int64(clip.ID),
		VideoID:   int64(clip.Videoid),
		URL:       clip.URL,
		Thumbnail: clip.Thumbnail,
	}
}
