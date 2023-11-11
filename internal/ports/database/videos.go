package database

import (
	"context"
	"database/sql"
	"funovation_23/internal/domain"
	"funovation_23/internal/ports/database/gen/funovation/public/model"
	"funovation_23/internal/ports/database/gen/funovation/public/table"
	"log"

	"github.com/go-jet/jet/v2/postgres"
)

type VideoStoreInterface interface {
	CreateVideo(ctx context.Context, video *domain.DownloadedYTVideo, videoPath, thumbnailPath string) (*domain.YoutubeVideo, error)
	GetVideo(ctx context.Context, id int64) (*domain.YoutubeVideo, error)
	GetAllVideos(ctx context.Context) ([]*domain.YoutubeVideo, error)
	CreateClip(ctx context.Context, clip domain.CreateClip) (*domain.Clip, error)
	CreateReel(ctx context.Context, reel domain.CreateReel) (*domain.Reel, error)
}

type VideosStore struct {
	db *sql.DB
}

func NewVideosStore(db *sql.DB) *VideosStore {
	return &VideosStore{db: db}
}

type VideoWithClips struct {
	Video model.Videos
	Clips []clip
}

type clip struct {
	model.Clips
}

func (vs *VideosStore) CreateVideo(ctx context.Context, video *domain.DownloadedYTVideo, videoPath, thumbnailPath string) (*domain.YoutubeVideo, error) {
	inssertModel := model.Videos{
		YoutubeID:   video.YoutubeID,
		Title:       video.Title,
		URL:         videoPath,
		Thumbnail:   thumbnailPath,
		CustomTitle: video.CustomTitle,
		Description: video.Description,
		S3IntPath:   video.IntS3Path,
	}
	stmt := table.Videos.INSERT(
		table.Videos.YoutubeID,
		table.Videos.Title,
		table.Videos.URL,
		table.Videos.CustomTitle,
		table.Videos.Thumbnail,
		table.Videos.Description,
		table.Videos.S3IntPath,
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
	stmt := postgres.SELECT(
		table.Videos.AllColumns,
		table.Clips.AllColumns,
	).
		FROM(table.Videos.INNER_JOIN(table.Clips, table.Videos.ID.EQ(table.Clips.Videoid))).
		WHERE(table.Videos.ID.EQ(postgres.Int(id))).
		GROUP_BY(table.Videos.ID, table.Clips.ID).
		ORDER_BY(table.Clips.Order.ASC())

	dest := []VideoWithClips{}
	err := stmt.QueryContext(ctx, vs.db, &dest)
	if err != nil {
		return nil, err
	}
	if len(dest) < 1 {
		return nil, domain.VideoNotFoundErr
	}
	mappedEntry := mapDBVideoWithClips(dest[0])
	log.Println(mappedEntry)
	return mappedEntry, nil
}

func (vs *VideosStore) GetAllVideos(ctx context.Context) ([]*domain.YoutubeVideo, error) {
	stmt := postgres.SELECT(
		table.Videos.AllColumns,
		table.Clips.AllColumns,
	).FROM(table.Videos.INNER_JOIN(table.Clips, table.Videos.ID.EQ(table.Clips.Videoid))).
		GROUP_BY(table.Videos.ID, table.Clips.ID).
		ORDER_BY(table.Videos.CreatedAt.DESC()).LIMIT(9)

	dest := []VideoWithClips{}
	err := stmt.QueryContext(ctx, vs.db, &dest)
	if err != nil {
		return nil, err
	}
	if len(dest) < 1 {
		return nil, nil
	}
	var result []*domain.YoutubeVideo
	for _, video := range dest {
		result = append(result, mapDBVideoWithClips(video))
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

func (vs *VideosStore) CreateReel(ctx context.Context, reel domain.CreateReel) (*domain.Reel, error) {
	insertModel := model.Reels{
		Videoid: int32(reel.VideoID),
		URL:     reel.URL,
	}
	stmt := table.Reels.INSERT(
		table.Reels.AllColumns.Except(table.Reels.ID, table.Reels.CreatedAt),
	).
		MODEL(insertModel).
		RETURNING(table.Reels.AllColumns)

	dest := model.Reels{}
	err := stmt.QueryContext(ctx, vs.db, &dest)
	if err != nil {
		return nil, err
	}
	return mapDBReel(dest), nil
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
		IntS3Path:   video.S3IntPath,
	}
}

func mapDBClip(clip model.Clips) *domain.Clip {
	return &domain.Clip{
		ID:        int64(clip.ID),
		VideoID:   int64(clip.Videoid),
		URL:       clip.URL,
		Thumbnail: clip.Thumbnail,
		Order:     int64(clip.Order),
	}
}

func mapDBVideoWithClips(videos VideoWithClips) *domain.YoutubeVideo {
	result := mapDBVideo(videos.Video)
	var clips []domain.Clip
	for _, clip := range videos.Clips {
		clips = append(clips, *mapDBClip(clip.Clips))
	}
	result.Clips = clips
	return result
}

func mapDBReel(reel model.Reels) *domain.Reel {
	return &domain.Reel{
		ID:      int64(reel.ID),
		VideoID: int64(reel.Videoid),
		URL:     reel.URL,
	}
}
