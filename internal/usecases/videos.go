package usecases

import (
	"context"
	"fmt"
	"funovation_23/internal/domain"
	"funovation_23/internal/ports/database"
	"funovation_23/internal/ports/downloader"
	fileStorage "funovation_23/internal/ports/filestore"
	"log"
)

type VideoUsecase struct {
	videoStore     database.VideoStoreInterface
	videoDowloader downloader.VideoDownloaderInterface
	storage        fileStorage.FileStorageInterface

	bucket string
}

func NewVideoUsecase(videoStore database.VideoStoreInterface, downloader downloader.VideoDownloaderInterface, storage fileStorage.FileStorageInterface, bucket string) *VideoUsecase {
	return &VideoUsecase{
		videoStore:     videoStore,
		videoDowloader: downloader,
		bucket:         bucket,
		storage:        storage,
	}
}

func (vu *VideoUsecase) ProcessYoutubeVideo(ctx context.Context, videoID string) (*domain.YoutubeVideo, error) {
	dowloadedVideo, err := vu.videoDowloader.DownloadYTVideo(videoID)
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("s3://%s%s", vu.bucket, dowloadedVideo.LocalPath)
	err = vu.storage.UploadFile(ctx, dowloadedVideo.LocalPath, path, "image/jpeg")
	if err != nil {
		return nil, err
	}
	log.Printf("File uploaded to %s\n", path)
	video, err := vu.videoStore.CreateVideo(ctx, dowloadedVideo, path)
	if err != nil {
		return nil, err
	}
	return video, nil
}

func (vu *VideoUsecase) GetAllVideos(ctx context.Context) ([]*domain.YoutubeVideo, error) {
	videos, err := vu.videoStore.GetAllVideos(ctx)
	if err != nil {
		return nil, err
	}
	return videos, nil
}

func (vu *VideoUsecase) GetVideo(ctx context.Context, id int64) (*domain.YoutubeVideo, error) {
	video, err := vu.videoStore.GetVideo(ctx, id)
	if err != nil {
		return nil, err
	}
	return video, nil
}
