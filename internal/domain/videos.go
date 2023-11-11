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
	IntS3Path   string
	Clips       []Clip
}

type DownloadedYTVideo struct {
	YoutubeID   string
	IntS3Path   string
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

type CreateClip struct {
	VideoID   int64
	URL       string
	Thumbnail string
	Order     int64
}

type Clip struct {
	ID        int64
	VideoID   int64
	URL       string
	Order     int64
	Thumbnail string
}

type Reel struct {
	ID      int64
	URL     string
	VideoID int64
}

type CreateReel struct {
	VideoID int64
	URL     string
}
