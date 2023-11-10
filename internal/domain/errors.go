package domain

import (
	"net/http"
)

type Error struct {
	Message string
	Code    int
}

func (e Error) Error() string { return e.Message }

var (
	VideoNotFoundErr = Error{Message: "Video was not found", Code: http.StatusNotFound}
)
