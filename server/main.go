package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

const (
	api  = "http://localhost:4080/api/mails/_search"
	user = "admin"
	pass = "Complexpass#123"
)

type Query struct {
	SearchType string         `json:"search_type"`
	Query_     map[string]any `json:"query"`
	SortFields []string       `json:"sort_fields"`
	From       int            `json:"from"`
	MaxResults int            `json:"max_results"`
	Highlight  map[string]any `json:"highlight"`
}

var defaultQuery = Query{
	SearchType: "match",
	Query_: map[string]any{
		"term":  "",
		"field": "Content",
	},
	SortFields: []string{},
	From:       0,
	MaxResults: 200,
	Highlight: map[string]any{
		"pre_tags":  []string{"<strong>"},
		"post_tags": []string{"</strong>"},
		"fields": map[string]any{
			"title": map[string]any{
				"pre_tags":  []string{},
				"post_tags": []string{},
			},
			"content": map[string]any{
				"pre_tags":  []string{},
				"post_tags": []string{},
			},
		},
	},
}

func main() {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Use(middleware.Logger)

	r.Get("/search", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("q")
		if q == "" {
			w.WriteHeader(400) // Bad Request
			w.Write([]byte("Search query parameter cannot be empty (q?)"))
		}

		defaultQuery.Query_["term"] = q

		payload, err := json.Marshal(defaultQuery)
		if err != nil {
			log.Fatal(err)
		}

		req, err := http.NewRequest("POST", api, bytes.NewBuffer(payload))
		if err != nil {
			log.Fatal(err)
		}
		req.SetBasicAuth(user, pass)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}

		if resp.StatusCode != 200 {
			w.WriteHeader(resp.StatusCode)
			w.Write([]byte(resp.Status))
			return
		}

		content, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		w.Write(content)
	})

	fmt.Println("HTTP Server listening and serving at :8080")
	http.ListenAndServe(":8080", r)
}
