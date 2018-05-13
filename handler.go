package main

import (
	"database/sql"
	"fmt"
	"net/http"
)

type Handler struct {
	db *sql.DB
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[1:]
	query := r.URL.RawQuery

	if path != "" {
		fmt.Fprintf(w, "%s", path)
	} else {
		fmt.Fprintf(w, "%s", query)
	}
}
