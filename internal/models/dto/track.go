package dto

type TrackResponse struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Artist string `json:"artist"`
	URL    string `json:"url"`
}
