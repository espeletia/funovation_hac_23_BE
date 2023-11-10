// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type Ack struct {
	Ok bool `json:"ok"`
}

type VideoRequest struct {
	URL         string `json:"url"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type VideoResponse struct {
	ID          string `json:"id"`
	URL         string `json:"url"`
	Title       string `json:"title"`
	YouTubeID   string `json:"youTubeId"`
	Status      string `json:"status"`
	Thumbnail   string `json:"thumbnail"`
	CustomTitle string `json:"customTitle"`
	Description string `json:"description"`
}
