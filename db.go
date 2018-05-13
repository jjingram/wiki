package main

import (
	"database/sql"

	"github.com/lib/pq"
)

type Page struct {
	Id    int
	Uri   string
	Title string
	Tags  []string
	Body  string
}

type DataAccessLayer interface {
	GetPage(string) (Page, error)
	PostPage(page Page) (sql.Result, error)
}

type DAL struct {
	db *sql.DB
}

func (dal DAL) GetPage(uri string) (Page, error) {
	var page Page
	err := dal.db.QueryRow("select id, uri, title, tags, body from page where uri = $1", uri).Scan(
		&page.Id, &page.Uri, &page.Title, pq.Array(&page.Tags), &page.Body)
	return page, err
}

func (dal DAL) PostPage(page Page) (sql.Result, error) {
	return dal.db.Exec("insert into page (uri, title, tags, body) values ($1, $2, $3, $4)",
		page.Uri, page.Title, pq.Array(page.Tags), page.Body)
}
