package api

import "time"

type Url struct {
	Url string `json:"url"`
}

type ErrorResponse struct {
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}
