package setup

import (
	"database/sql"
	"funovation_23/graph"
	"funovation_23/internal/config"
	"funovation_23/internal/ports/database"
	"funovation_23/internal/ports/downloader"
	fileStorage "funovation_23/internal/ports/filestore"
	"funovation_23/internal/usecases"
	"funovation_23/internal/usecases/encoding/images"
	"funovation_23/internal/usecases/encoding/video"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func NewResolver(dbConn *sql.DB, config config.Config, s3Client *s3.Client) (*graph.Resolver, error) {
	videoStore := database.NewVideosStore(dbConn)
	VideoDownloader := downloader.NewVideoDownloader()
	fileStorage := fileStorage.NewFileS3Storage(s3Client)

	imageEncoding := images.NewImageMediaEncoder(config.EncodingConfig.FfmpegPath)
	videoEncoding := video.NewVideoMediaEncoder(config.EncodingConfig.FfmpegPath, config.EncodingConfig.FfprobePath, config.EncodingConfig.MaxClipLength)

	videoUsecase := usecases.NewVideoUsecase(videoStore, VideoDownloader, fileStorage, config.S3Config.Bucket, imageEncoding, videoEncoding, config.S3Config.Prod)
	return &graph.Resolver{
		VideoUsecase: videoUsecase,
		Mapper:       graph.NewMapper(),
		InputMapper:  graph.NewInputMapper(),
	}, nil

}
