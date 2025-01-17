package models

type User struct {
	ID       int
	Login    string
	Email    string
	Password string
	Sale     string
}

type UserRequest struct {
	ID       int    `json:"id"`
	Login    string `json:"login"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID    int    `json:"id"`
	Login string `json:"login"`
	Email string `json:"email"`
}
