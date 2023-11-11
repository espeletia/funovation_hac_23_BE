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
		Clips:       m.mapClips(ytVideo.Clips),
	}
}

func (m *Mapper) mapClips(clips []domain.Clip) []*model.Clip {
	gqlClips := []*model.Clip{}
	for _, clip := range clips {
		gqlClips = append(gqlClips, &model.Clip{
			ID:        strconv.FormatInt(clip.ID, 10),
			VideoID:   strconv.FormatInt(clip.VideoID, 10),
			URL:       clip.URL,
			Thumbnail: clip.Thumbnail,
		})
	}
	return gqlClips
}

func (m *Mapper) mapReel(reel *domain.Reel) *model.Reel {
	return &model.Reel{
		ID:      strconv.FormatInt(reel.ID, 10),
		VideoID: strconv.FormatInt(reel.VideoID, 10),
		URL:     reel.URL,
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
