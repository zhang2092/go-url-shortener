package handler

import (
	"io/fs"
	"net/http"

	"github.com/zhang2092/go-url-shortener/db"
)

func HomeView(templates fs.FS, store db.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		user := withUser(ctx)
		result, err := store.ListUrlByUser(ctx, user.ID)
		if err != nil {
			renderLayout(w, r, templates, nil, "home.html.tmpl")
		}

		scheme := "http://"
		if r.TLS != nil {
			scheme = "https://"
		}
		for _, item := range result {
			item.ShortUrl = scheme + r.Host + "/" + item.ShortUrl
		}
		renderLayout(w, r, templates, result, "home.html.tmpl")
	}
}
