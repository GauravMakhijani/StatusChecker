package db

import (
	"github.com/jmoiron/sqlx"
)

type WebsiteStatus struct {
	ID     int64  `db:"id"`
	Link   string `db:"link"`
	Status string `db:"status"`
}

// StatusStorer is an interface for the database
type StatusStorer interface {
	InsertWebsite(website WebsiteStatus) (err error)
	UpdateWebsiteStatus(url string, status string) (err error)
	GetWebsiteStatus(name string) (ws []WebsiteStatus, err error)
	GetAll() (ws []WebsiteStatus, err error)
}

// store is a struct that implements the StatusStorer interface
type store struct {
	DB *sqlx.DB
}

// New returns a new store
func New(db *sqlx.DB) StatusStorer {
	return &store{DB: db}
}

// UpdateWebsiteStatus updates the status of a website
func (s *store) UpdateWebsiteStatus(url string, status string) (err error) {

	// fmt.Println("Checking status of", ws.Link)
	s.DB.Exec("UPDATE links SET status = $2 where link = $1", url, status)

	return
}

// GetWebsiteStatus returns a particular website
func (s *store) GetWebsiteStatus(query string) (ws []WebsiteStatus, err error) {

	ws = []WebsiteStatus{}
	err = s.DB.Select(&ws, "SELECT id,link,status FROM links WHERE link LIKE $1", "%"+query+"%")

	return
}

// GetAll returns all the websites in the database
func (s *store) GetAll() (ws []WebsiteStatus, err error) {

	ws = []WebsiteStatus{}
	err = s.DB.Select(&ws, "SELECT id,link,status FROM links")

	return
}

// InsertWebsite inserts a website into the database
func (s *store) InsertWebsite(website WebsiteStatus) (err error) {

	s.DB.Exec("INSERT INTO links (link, status) VALUES ($1, $2)", website.Link, website.Status)

	return
}

// ROuter -> handler -> service -> repo -> DB{mysql, postgres, mongo, redis}
