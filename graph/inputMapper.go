package graph

import (
	"funovation_23/graph/model"
	"funovation_23/internal/domain"
	"funovation_23/internal/util"
)

type InputMapper struct{}

func NewInputMapper() *InputMapper {
	return &InputMapper{}
}

func (m *InputMapper) MapCreateVideoInputToDomain(input model.VideoRequest) domain.CreateVideo {
	videoID := util.GetYoutubeIDFromURL(input.URL)

	return domain.CreateVideo{
		YoutubeID:   videoID,
		CustomTitle: input.Title,
		Description: input.Description,
	}
}
