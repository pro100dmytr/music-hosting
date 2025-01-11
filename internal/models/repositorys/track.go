package repositorys

// TODO: move to repository/models.go
type Track struct {
	ID       int
	Name     string
	Artist   string
	URL      string
	Likes    int
	Dislikes int
}
