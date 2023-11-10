package graph

import (
	"funovation_23/graph/model"
	"funovation_23/internal/domain"
	"strconv"
)

type Mapper struct{}

func NewMapper() *Mapper {
	return &Mapper{}
}

func (m *Mapper) mapYoutubeVideoToGqlVideo(ytVideo *domain.YoutubeVideo) *model.VideoResponse {
	return &model.VideoResponse{
		ID:          strconv.FormatInt(ytVideo.ID, 10),
		Title:       ytVideo.Title,
		URL:         ytVideo.S3Path,
		YouTubeID:   ytVideo.YoutubeID,
		Status:      m.mapStatus(ytVideo.Status),
		Thumbnail:   ytVideo.Thumbnail,
		Description: ytVideo.Description,
		CustomTitle: ytVideo.CustomTitle,
	}
}

func (m *Mapper) mapStatus(status int64) string {
	switch status {
	case 0:
		return "pending"
	case 1:
		return "processing"
	case 2:
		return "done"
	default:
		return "unknown"
	}
}
