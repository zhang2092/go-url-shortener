package handler

import (
	"io/fs"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/zhang2092/go-url-shortener/db"
	"github.com/zhang2092/go-url-shortener/service"
	"github.com/zhang2092/go-url-shortener/shortener"
)

func CreateShortUrlView(templates fs.FS) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		renderLayout(w, r, templates, nil, "short_url/create.html.tmpl")
	}
}

func CreateShortUrl(templates fs.FS, store db.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		if err := r.ParseForm(); err != nil {
			renderCreateShortUrl(w, r, templates, map[string]string{"Error": "请求参数错误"})
			return
		}

		ctx := r.Context()
		user := withUser(ctx)
		longUrl := r.PostFormValue("long_url")
		shortUrl, err := shortener.GenerateShortLink(longUrl, user.ID)
		if err != nil {
			renderCreateShortUrl(w, r, templates, map[string]string{"Error": "生成短路径错误"})
			return
		}

		log.Println(shortUrl)

		_, err = store.CreateUserUrl(ctx, &db.CreateUserUrlParams{
			UserID:    user.ID,
			ShortUrl:  shortUrl,
			OriginUrl: longUrl,
			ExpireAt:  time.Now().Add(time.Hour * 6),
		})
		if err != nil {
			renderCreateShortUrl(w, r, templates, map[string]string{"Error": "短路径存储错误"})
			return
		}

		err = service.SaveUrlMapping(shortUrl, longUrl, user.ID)
		if err != nil {
			renderCreateShortUrl(w, r, templates, map[string]string{"Error": "短路径存储错误"})
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func HandleShortUrlRedirect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shorUrl := vars["shortUrl"]
	link, err := service.RetrieveInitialUrl(shorUrl)
	if err != nil {
		http.Error(w, "failed to get url", http.StatusInternalServerError)
		return
	}
	if len(link) == 0 {
		http.Error(w, "short url get to long url is empty", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, link, http.StatusFound)
}

func renderCreateShortUrl(w http.ResponseWriter, r *http.Request, templates fs.FS, data any) {
	renderLayout(w, r, templates, data, "short_url/create.html.tmpl")
}
