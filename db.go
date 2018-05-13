package main

import (
	"database/sql"
)

type Page struct {
	Id    int
	Title []byte
	Tags  [][]byte
	Body  []byte
}

type DataAccessLayer interface {
	GetPageById(int) (Page, error)
	GetPageByTitle(string) (Page, error)
}

type DAL struct {
	db *sql.DB
}

func (dal DAL) GetPageById(id int) (Page, error) {
	var page Page
	err := dal.db.QueryRow("select id, title, tags, body from page where id = ?", id).Scan(&page)
	return page, err
}

func (dal DAL) GetPageByTitle(title string) (Page, error) {
	var page Page
	err := dal.db.QueryRow("select id, title, tags, body from page where title = '?'", title).Scan(&page)
	return page, err
}
