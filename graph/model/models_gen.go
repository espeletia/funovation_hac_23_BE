// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type Ack struct {
	Ok bool `json:"ok"`
}

type VideoResponse struct {
	ID        string `json:"id"`
	URL       string `json:"url"`
	Title     string `json:"title"`
	Duration  int    `json:"duration"`
	Thumbnail string `json:"thumbnail"`
}
