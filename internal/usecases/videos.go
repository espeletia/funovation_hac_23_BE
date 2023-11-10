package usecases

import (
	"context"
	"fmt"
	"funovation_23/internal/domain"
	"funovation_23/internal/ports/database"
	"funovation_23/internal/ports/downloader"
	fileStorage "funovation_23/internal/ports/filestore"
	"funovation_23/internal/usecases/encoding/images"
	"log"
	"strings"
)

type VideoUsecase struct {
	videoStore     database.VideoStoreInterface
	videoDowloader downloader.VideoDownloaderInterface
	storage        fileStorage.FileStorageInterface

	imageEncoder *images.ImageMediaEncoder

	bucket string
}

func NewVideoUsecase(videoStore database.VideoStoreInterface, downloader downloader.VideoDownloaderInterface, storage fileStorage.FileStorageInterface, bucket string, imageEncoder *images.ImageMediaEncoder) *VideoUsecase {
	return &VideoUsecase{
		videoStore:     videoStore,
		videoDowloader: downloader,
		bucket:         bucket,
		storage:        storage,
		imageEncoder:   imageEncoder,
	}
}

func (vu *VideoUsecase) ProcessYoutubeVideo(ctx context.Context, videoCreate domain.CreateVideo) (*domain.YoutubeVideo, error) {
	dowloadedVideo, err := vu.videoDowloader.DownloadYTVideo(videoCreate)
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("s3://%s/video%s", vu.bucket, dowloadedVideo.LocalPath)
	uploadedVideoPath, err := vu.storage.UploadFile(ctx, dowloadedVideo.LocalPath, path, "video/mp4")
	if err != nil {
		return nil, err
	}
	imagePath := strings.Replace(dowloadedVideo.LocalPath, "video", "image", 1)
	imagePath = strings.Replace(imagePath, ".mp4", ".jpg", 1)
	err = vu.imageEncoder.GenerateThumbanail(ctx, dowloadedVideo.LocalPath, imagePath)
	if err != nil {
		return nil, err
	}
	path = fmt.Sprintf("s3://%s/thumbnail%s", vu.bucket, imagePath)
	uploadedThumbnailPath, err := vu.storage.UploadFile(ctx, imagePath, path, "image/jpeg")
	if err != nil {
		return nil, err
	}
	log.Printf("File uploaded to %s\n", path)
	video, err := vu.videoStore.CreateVideo(ctx, dowloadedVideo, uploadedVideoPath, uploadedThumbnailPath)
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
