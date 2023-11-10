package domain

type YoutubeVideo struct {
	ID        int64
	Title     string
	S3Path    string
	YoutubeID string
	Status    int64
}

type DownloadedYTVideo struct {
	YoutubeID string
	LocalPath string
	Title     string
}
