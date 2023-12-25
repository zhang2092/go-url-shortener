package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/zhang2092/go-url-shortener/shortener"
	"github.com/zhang2092/go-url-shortener/store"
)

type UrlCreationRequest struct {
	LongUrl string `json:"long_url"`
	UserId  string `json:"user_id"`
}

type UrlCreationResponse struct {
	Message  string `json:"message"`
	ShortUrl string `json:"short_url"`
}

func CreateShortUrl(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req UrlCreationRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid parameter", http.StatusInternalServerError)
		return
	}

	shortUrl, err := shortener.GenerateShortLink(req.LongUrl, req.UserId)
	if err != nil {
		http.Error(w, "failed to generate short link", http.StatusInternalServerError)
		return
	}

	err = store.SaveUrlMapping(shortUrl, req.LongUrl, req.UserId)
	if err != nil {
		http.Error(w, "failed to store url mapping", http.StatusInternalServerError)
		return
	}

	scheme := "http://"
	if r.TLS != nil {
		scheme = "https://"
	}

	res := &UrlCreationResponse{
		Message:  "short url created successfully",
		ShortUrl: scheme + r.Host + "/" + shortUrl,
	}

	b, err := json.Marshal(res)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func HandleShortUrlRedirect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shorUrl := vars["shortUrl"]
	link, err := store.RetrieveInitialUrl(shorUrl)
	if err != nil {
		http.Error(w, "failed to get url", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, link, http.StatusFound)
}
