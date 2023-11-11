package usecases

import (
	"context"
	"fmt"
	"funovation_23/internal/domain"
	"funovation_23/internal/ports/database"
	"funovation_23/internal/ports/downloader"
	fileStorage "funovation_23/internal/ports/filestore"
	"funovation_23/internal/usecases/encoding/images"
	"funovation_23/internal/usecases/encoding/video"
	"log"
	"os"
	"strings"
)

type VideoUsecase struct {
	videoStore     database.VideoStoreInterface
	videoDowloader downloader.VideoDownloaderInterface
	storage        fileStorage.FileStorageInterface

	imageEncoder *images.ImageMediaEncoder
	videoEncoder *video.VideoMediaEncoder

	bucket string
	prod   bool
}

func NewVideoUsecase(videoStore database.VideoStoreInterface, downloader downloader.VideoDownloaderInterface, storage fileStorage.FileStorageInterface, bucket string, imageEncoder *images.ImageMediaEncoder, videoEncoder *video.VideoMediaEncoder, prod bool) *VideoUsecase {
	return &VideoUsecase{
		videoStore:     videoStore,
		videoDowloader: downloader,
		bucket:         bucket,
		storage:        storage,
		imageEncoder:   imageEncoder,
		videoEncoder:   videoEncoder,
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
	dowloadedVideo.IntS3Path = path
	if vu.prod {
		dowloadedVideo.IntS3Path = strings.Replace(uploadedVideoPath, "https://funovation.fra1.digitaloceanspaces.com/", "s3://funovation/", 1)
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
	err = os.Remove(dowloadedVideo.LocalPath)
	if err != nil {
		return nil, err
	}
	err = os.Remove(imagePath)
	if err != nil {
		return nil, err
	}
	log.Printf("File uploaded to %s\n", path)
	video, err := vu.videoStore.CreateVideo(ctx, dowloadedVideo, uploadedVideoPath, uploadedThumbnailPath)
	if err != nil {
		return nil, err
	}
	clips, err := vu.CreateClips(ctx, *video)
	if err != nil {
		return nil, err
	}
	video.Clips = clips
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

func (vu *VideoUsecase) CreateClips(ctx context.Context, video domain.YoutubeVideo) ([]domain.Clip, error) {
	path := video.IntS3Path
	videoFile, err := vu.storage.DownloadFile(ctx, path, "")
	if err != nil {
		log.Println("HAIIIIII")
		return nil, err
	}
	defer func() {
		err = videoFile.Close()
		if err != nil {
			log.Println(err, "FUCKKK!!!!!")
		}
	}()
	_, err = vu.videoEncoder.GenerateClips(ctx, videoFile.Name(), "path")
	if err != nil {
		log.Println("ENCODER FAIL >_<:", err)
		return nil, err
	}
	files, err := os.ReadDir("path")
	if err != nil {
		log.Println("READDIR() FAIL >_<:", err)
		return nil, err
	}
	var clips []domain.Clip
	for i, file := range files {
		if file.IsDir() {
			continue
		}

		// path := "s3://" + vu.bucket + "/clips/" + file.Name()
		path := fmt.Sprintf("s3://%s/clips/%s", vu.bucket, file.Name())
		uploadFile, err := vu.storage.UploadFile(ctx, "path/"+file.Name(), path, "video/mp4")
		if err != nil {
			return nil, err
		}
		imagePath := strings.Replace(file.Name(), "clip", "image", 1)
		imagePath = strings.Replace(imagePath, ".mp4", ".jpg", 1)
		// path := fmt.Sprintf("s3://%s/image/%s", vu.bucket)
		err = vu.imageEncoder.GenerateThumbanail(ctx, "path/"+file.Name(), imagePath)
		if err == domain.UnableToThumbnail {
			clip, err := vu.videoStore.CreateClip(ctx, domain.CreateClip{
				VideoID:   video.ID,
				URL:       uploadFile,
				Thumbnail: "",
				Order:     int64(i),
			})
			if err != nil {
				return nil, err
			}
			clips = append(clips, *clip)
			os.Remove("path/" + file.Name())
			os.Remove("path/" + imagePath)
			continue
		}
		if err != nil {
			return nil, err
		}
		path = fmt.Sprintf("s3://%s/thumbnail/%s", vu.bucket, imagePath)
		uploadImageFile, err := vu.storage.UploadFile(ctx, "path/"+file.Name(), path, "image/jpeg")
		if err != nil {
			return nil, err
		}
		clip, err := vu.videoStore.CreateClip(ctx, domain.CreateClip{
			VideoID:   video.ID,
			URL:       uploadFile,
			Thumbnail: uploadImageFile,
			Order:     int64(i),
		})
		if err != nil {
			return nil, err
		}
		clips = append(clips, *clip)
		os.Remove("path/" + file.Name())
		os.Remove("path/" + imagePath)
	}

	return clips, nil
}
