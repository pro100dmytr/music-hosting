package repositorys

type User struct {
	ID         int
	Login      string
	Email      string
	Password   string
	PlaylistID []int
	Salt       string
}
