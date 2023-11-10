package downloader

import (
	"funovation_23/internal/domain"
	"io"
	"os"
	"path/filepath"

	"github.com/kkdai/youtube/v2"
)

type VideoDownloaderInterface interface {
	DownloadYTVideo(videoCreate domain.CreateVideo) (*domain.DownloadedYTVideo, error)
}

type VideoDownloader struct{}

func NewVideoDownloader() *VideoDownloader { return &VideoDownloader{} }

func (YTD *VideoDownloader) DownloadYTVideo(videoCreate domain.CreateVideo) (*domain.DownloadedYTVideo, error) {
	client := youtube.Client{}

	video, err := client.GetVideo(videoCreate.YoutubeID)
	if err != nil {
		return nil, err
	}

	formats := video.Formats.WithAudioChannels() // only get videos with audio
	stream, _, err := client.GetStream(video, &formats[0])
	if err != nil {
		return nil, err
	}
	defer stream.Close()

	file, err := os.CreateTemp("", "video*.mp4")
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(file, stream)
	if err != nil {
		return nil, err
	}

	path, err := filepath.Abs(file.Name())
	if err != nil {
		return nil, err
	}
	// pathSplice := filepath.SplitList(path)
	// lastIndex := len(pathSplice) - 1
	// path = pathSplice[lastIndex]
	return &domain.DownloadedYTVideo{
		YoutubeID:   videoCreate.YoutubeID,
		LocalPath:   path,
		Title:       video.Title,
		CustomTitle: videoCreate.CustomTitle,
		Description: videoCreate.Description,
	}, nil
}
