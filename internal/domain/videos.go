package domain

type YoutubeVideo struct {
	ID        int64
	Title     string
	S3Path    string
	YoutubeID string
	LocalPath string
}

type DownloadedYTVideo struct {
	YoutubeID string
	LocalPath string
	Title     string
}
