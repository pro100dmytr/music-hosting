package dto

type UserResponse struct {
	ID          int    `json:"id"`
	Login       string `json:"login"`
	Email       string `json:"email"`
	PlaylistsID []int  `json:"playlist_id"`
}
