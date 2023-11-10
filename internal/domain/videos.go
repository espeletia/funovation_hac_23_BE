package domain

type YoutubeVideo struct {
	ID          int64
	Title       string
	S3Path      string
	YoutubeID   string
	Status      int64
	Thumbnail   string
	CustomTitle string
	Description string
}

type DownloadedYTVideo struct {
	YoutubeID   string
	CustomTitle string
	Description string
	LocalPath   string
	Title       string
}

type CreateVideo struct {
	YoutubeID   string
	CustomTitle string
	Description string
}
