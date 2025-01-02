package https

type Track struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Artist   string `json:"artist"`
	URL      string `json:"url"`
	Likes    int    `json:"likes"`
	Dislikes int    `json:"dislikes"`
}
