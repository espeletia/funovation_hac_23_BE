package graph

import (
	"funovation_23/graph/model"
	"funovation_23/internal/domain"
	"funovation_23/internal/util"
	"strconv"
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

func (m *InputMapper) MapCreateReel(input model.ReelRequest) (int64, []int64, error) {
	clipIDs := []int64{}
	for _, id := range input.ClipIds {
		parsedId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return 0, nil, err
		}
		clipIDs = append(clipIDs, parsedId)
	}
	parsedVideoID, err := strconv.ParseInt(input.VideoID, 10, 64)
	if err != nil {
		return 0, nil, err
	}
	return parsedVideoID, clipIDs, nil
}
