package main

import (
	"database/sql"
)

type Page struct {
	Id    int
	Uri   string
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

func (dal DAL) GetPage(uri string) (Page, error) {
	var page Page
	err := dal.db.QueryRow("select id, uri, title, tags, body from page where uri = $1",
		title).Scan(&page.Id, &page.Uri, &page.Title, &page.Tags, &page.Body)
	return page, err
}
