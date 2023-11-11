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
	"sort"
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
		prod:           prod,
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
	path := video.S3Path
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

func (vu *VideoUsecase) CreateReel(ctx context.Context, id int64, clipIDs []int64) (*domain.Reel, error) {
	reel, videoId, err := vu.createReel(ctx, id, clipIDs)
	if err != nil {
		return nil, err
	}
	uploadPath := fmt.Sprintf("s3://%s/reel%s", vu.bucket, reel)
	log.Println("UPLOADPATH:", uploadPath)
	log.Print("VIDEOID:", reel)
	uploadedReel, err := vu.storage.UploadFile(ctx, reel, uploadPath, "video/mp4")
	if err != nil {
		return nil, err
	}
	storedReel, err := vu.videoStore.CreateReel(ctx, domain.CreateReel{
		VideoID: videoId,
		URL:     uploadedReel,
	})
	if err != nil {
		return nil, err
	}
	return storedReel, nil
}

func (vu *VideoUsecase) createReel(ctx context.Context, id int64, clipIDs []int64) (string, int64, error) {
	video, err := vu.videoStore.GetVideo(ctx, id)
	if err != nil {
		return "", 0, err
	}

	tmpPath := fmt.Sprintf("tmp/%d", video.ID)
	err = os.MkdirAll(tmpPath, 0755)
	if err != nil {
		return "", 0, err
	}
	var videoClipIDs []int64
	for _, clip := range video.Clips {
		videoClipIDs = append(videoClipIDs, clip.ID)
	}
	log.Println("ORIGINALIDS:", videoClipIDs)
	log.Println("REQUEST:", clipIDs)
	// Check if the sets of clip IDs are equal
	IDsMap := map[int64]domain.Clip{} // create and initialize a map of IDs from the new queue
	for _, entry := range video.Clips {
		IDsMap[entry.ID] = entry
	}
	var clips []domain.Clip // create a new queue that will be stored in the database
	for _, entry := range clipIDs {
		value, ok := IDsMap[entry]
		if !ok { // verify that the new queue contains all the IDs from the current queue
			return "", 0, fmt.Errorf("invalid clip IDs")
		}
		clips = append(clips, value)
	}

	// If the sets are equal, use the clips from the current queue

	sort.Slice(clips, func(i, j int) bool {
		return clips[i].Order < clips[j].Order
	})
	for _, clip := range clips {
		vu.storage.DownloadFile(ctx, clip.URL, tmpPath)
	}
	files, err := os.ReadDir(tmpPath)
	if err != nil {
		return "", 0, err
	}
	filePaths := []string{}
	for _, file := range files {
		filePaths = append(filePaths, tmpPath+"/"+file.Name())
		defer os.Remove(tmpPath + "/" + file.Name())
	}
	log.Println(filePaths)
	resultFile, err := vu.videoEncoder.CreateReel(ctx, filePaths)
	if err != nil {
		return "", 0, err
	}
	return resultFile, video.ID, nil
}
