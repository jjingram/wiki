package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
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

		title := bluemonday.UGCPolicy().SanitizeBytes([]byte(page.Title))

		unsafe := blackfriday.Run([]byte(page.Body))
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
			fmt.Fprint(w, http.StatusText(http.StatusInternalServerError))
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

	headers := make(map[string][]string)
	scanner := bufio.NewScanner(body)
	for scanner.Scan() {
		txt := scanner.Text()

		if txt == "" {
			break
		}

		header := strings.Split(txt, ":")

		if len(header) != 2 {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, http.StatusText(http.StatusBadRequest))
			return
		}

		key := strings.ToLower(strings.Trim(header[0], " \t"))
		val := strings.ToLower(strings.Trim(header[1], " \t"))

		_, ok := headers[key]
		if !ok {
			headers[key] = make([]string, 0)
		}
		headers[key] = append(headers[key], val)
	}

	if headers["title"] == nil || headers["tag"] == nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, http.StatusText(http.StatusBadRequest))
		return
	}

	var (
		page Page
		err  error
	)
	page.Title = headers["title"][0]
	page.Uri = hyphenate(page.Title)
	page.Tags = headers["tag"]

	for scanner.Scan() {
		page.Body += scanner.Text() + "\n"
	}

	_, err = h.dal.PostPage(page)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, http.StatusText(http.StatusInternalServerError))
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, http.StatusText(http.StatusCreated))
}

func hyphenate(text string) string {
	var hyphenated []rune
	var hyphen = false
	for _, r := range text {
		switch {
		case unicode.IsLetter(r) || unicode.IsNumber(r):
			if hyphen && len(hyphenated) > 0 {
				hyphenated = append(hyphenated, '-')
			}
			hyphen = false
			hyphenated = append(hyphenated, unicode.ToLower(r))
		default:
			hyphen = true
		}
	}
	return string(hyphenated)
}
