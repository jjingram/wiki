package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/microcosm-cc/bluemonday"
	"gopkg.in/russross/blackfriday.v2"
)

type Handler struct {
	dal DataAccessLayer
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
		page, err := h.dal.GetPage(path)
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, http.StatusText(http.StatusNotFound))
		} else if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, http.StatusText(http.StatusInternalServerError))
		}

		unsafe := blackfriday.Run(page.Body)
		html := bluemonday.UGCPolicy().SanitizeBytes(unsafe)
		fmt.Fprint(w, string(html))
	} else if query != "" {
		fmt.Fprint(w, query)
		return
	}
}

func (h *Handler) doPost(w http.ResponseWriter, r *http.Request) {
	body := r.Body
	defer body.Close()

	input, err := ioutil.ReadAll(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, http.StatusText(http.StatusInternalServerError))
		return
	}
	fmt.Fprint(w, input)
}
