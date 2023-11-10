package usecases

import (
	"context"
	"fmt"
	"funovation_23/internal/ports/database"
	"funovation_23/internal/ports/downloader"
	fileStorage "funovation_23/internal/ports/filestore"
	"log"
)

type VideoUsecase struct {
	videoStore     database.VideoStoreInterface
	videoDowloader downloader.VideoDownloaderInterface
	storage        fileStorage.FileStorageInterface
}

func NewVideoUsecase(videoStore database.VideoStoreInterface, downloader downloader.VideoDownloaderInterface, storage fileStorage.FileStorageInterface) *VideoUsecase {
	return &VideoUsecase{
		videoStore:     videoStore,
		videoDowloader: downloader,
		storage:        storage,
	}
}

func (vu *VideoUsecase) ProcessYoutubeVideo(videoID string) error {
	dowloadedVideo, err := vu.videoDowloader.DownloadYTVideo(videoID)
	if err != nil {
		return err
	}
	log.Println("Downloaded video", dowloadedVideo)
	err = vu.storage.UploadFile(context.Background(), dowloadedVideo.LocalPath, fmt.Sprintf("s3://test/uploads%s", dowloadedVideo.LocalPath), "image/jpeg")
	if err != nil {
		return err
	}

	return nil
}
