package repositorys

// TODO: move to repository/models.go
type User struct {
	ID         int
	Login      string
	Email      string
	Password   string
	PlaylistID []int
	Salt       string
}
