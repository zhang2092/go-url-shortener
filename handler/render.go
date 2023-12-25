package handler

import (
	"html/template"
	"io/fs"
	"net/http"
	"path/filepath"

	"github.com/gorilla/csrf"
	"github.com/zhang2092/go-url-shortener/pkg/logger"
)

// renderLayout 渲染方法 带框架
func renderLayout(w http.ResponseWriter, r *http.Request, templates fs.FS, data any, tmpl string) {
	t := template.New(filepath.Base(tmpl))
	t = t.Funcs(template.FuncMap{
		"csrfField": func() template.HTML {
			return csrf.TemplateField(r)
		},
		"currentUser": func() *Authorize {
			return withUser(r.Context())
		},
		"genShortUrl": func(url string) string {
			scheme := "http://"
			if r.TLS != nil {
				scheme = "https://"
			}
			return scheme + r.Host + "/" + url
		},
	})

	tpl := template.Must(t.Clone())
	tpl, err := tpl.ParseFS(templates, tmpl, "base/header.html.tmpl", "base/footer.html.tmpl")
	if err != nil {
		logger.Logger.Errorf("template parse: %s, %v", tmpl, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := tpl.Execute(w, data); err != nil {
		logger.Logger.Errorf("template execute: %s, %v", tmpl, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
