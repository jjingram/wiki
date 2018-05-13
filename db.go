package main

import (
	"database/sql"
)

type Page struct {
	Id    int
	Title []byte
	Tags  StringSlice
	Body  []byte
}

type DataAccessLayer interface {
	GetPage(string) (Page, error)
}

type DAL struct {
	db *sql.DB
}

func (dal DAL) GetPage(title string) (Page, error) {
	var page Page
	err := dal.db.QueryRow("select id, title, tags, body from page where title = $1",
		title).Scan(&page.Id, &page.Title, &page.Tags, &page.Body)
	return page, err
}
