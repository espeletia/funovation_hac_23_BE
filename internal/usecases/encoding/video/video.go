package video

type VideoMediaEncoder struct {
	ffmpegPath string
}

func NewVideoMediaEncoder(ffmpegPath string) *VideoMediaEncoder {
	return &VideoMediaEncoder{ffmpegPath: ffmpegPath}
}
