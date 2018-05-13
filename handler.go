package main

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
)

type Handler struct {
	db *sql.DB
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.doGet(w, r)
	case http.MethodPost:
		h.doPost(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, http.StatusText(http.StatusMethodNotAllowed))
	}
}

func (h *Handler) doGet(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[1:]
	query := r.URL.RawQuery

	if path != "" {
		fmt.Fprintf(w, "%s", path)
	} else {
		fmt.Fprintf(w, "%s", query)
	}
}

func (h *Handler) doPost(w http.ResponseWriter, r *http.Request) {
	body := r.Body
	defer body.Close()

	_, err := io.Copy(w, body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%s", http.StatusText(http.StatusInternalServerError))
		return
	}
}
