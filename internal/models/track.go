package models

type Track struct {
	ID       int
	Name     string
	Artist   string
	URL      string
	Likes    int
	Dislikes int
}

type TrackRequest struct {
	Name     string `json:"name"`
	Artist   string `json:"artist"`
	URL      string `json:"url"`
	Likes    int    `json:"likes"`
	Dislikes int    `json:"dislikes"`
}

type TrackResponse struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Artist   string `json:"artist"`
	URL      string `json:"url"`
	Likes    int    `json:"likes"`
	Dislikes int    `json:"dislikes"`
}
