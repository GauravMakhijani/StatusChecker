package db

type Website struct {
	Link   string `db:"link"`
	Status string `db:"status"`
}
