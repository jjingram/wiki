package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"unicode"

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
			return
		} else if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, http.StatusText(http.StatusInternalServerError))
			return
		}

		title := bluemonday.UGCPolicy().SanitizeBytes(page.Title)

		unsafe := blackfriday.Run(page.Body)
		body := bluemonday.UGCPolicy().SanitizeBytes(unsafe)

		t, err := template.New("layout").ParseFiles("layout.html")
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, http.StatusText(http.StatusInternalServerError))
			return
		}
		data := struct {
			Title string
			Body  template.HTML
		}{
			Title: string(title),
			Body:  template.HTML(body),
		}
		err = t.Execute(w, data)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, http.StatusText(http.StatusInternalServerError))
			return
		}
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

func hyphenate(text string) string {
	var anchorName []rune
	var hyphen = false
	for _, r := range text {
		switch {
		case unicode.IsLetter(r) || unicode.IsNumber(r):
			if hyphen && len(anchorName) > 0 {
				anchorName = append(anchorName, '-')
			}
			hyphen = false
			anchorName = append(anchorName, unicode.ToLower(r))
		default:
			hyphen = true
		}
	}
	return string(anchorName)
}
