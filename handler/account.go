package handler

import (
	"database/sql"
	"io/fs"
	"net/http"
	"time"

	"github.com/zhang2092/go-url-shortener/db"
	"github.com/zhang2092/go-url-shortener/pkg/cookie"
	pwd "github.com/zhang2092/go-url-shortener/pkg/password"
)

type registerPageData struct {
	Summary     string
	Email       string
	EmailMsg    string
	Username    string
	UsernameMsg string
	Password    string
	PasswordMsg string
}

func RegisterView(templates fs.FS) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		renderRegister(w, r, templates, nil)
	}
}

func Register(templates fs.FS, store db.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		email := r.PostFormValue("email")
		username := r.PostFormValue("username")
		password := r.PostFormValue("password")
		resp, ok := viladatorRegister(email, username, password)
		if !ok {
			renderRegister(w, r, templates, resp)
			return
		}

		hashedPassword, err := pwd.BcryptHashPassword(password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		arg := &db.CreateUserParams{
			ID:             genId(),
			Username:       username,
			HashedPassword: hashedPassword,
			Email:          email,
		}

		_, err = store.CreateUser(r.Context(), arg)
		if err != nil {
			if store.IsUniqueViolation(err) {
				resp.Summary = "邮箱或名称已经存在"
				renderRegister(w, r, templates, resp)
				return
			}

			resp.Summary = "请求网络错误,请刷新重试"
			renderRegister(w, r, templates, resp)
			return
		}

		http.Redirect(w, r, "/login", http.StatusFound)
	}
}

type loginPageData struct {
	Summary     string
	Email       string
	EmailMsg    string
	Password    string
	PasswordMsg string
}

func LoginView(templates fs.FS) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		renderLogin(w, r, templates, nil)
	}
}

func Login(templates fs.FS, store db.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		if err := r.ParseForm(); err != nil {
			renderLogin(w, r, templates, registerPageData{Summary: "请求网络错误,请刷新重试"})
			return
		}

		email := r.PostFormValue("email")
		password := r.PostFormValue("password")
		resp, ok := viladatorLogin(email, password)
		if !ok {
			renderLogin(w, r, templates, resp)
			return
		}

		ctx := r.Context()
		user, err := store.GetUserByEmail(ctx, email)
		if err != nil {
			if store.IsNoRows(sql.ErrNoRows) {
				resp.Summary = "邮箱或密码错误"
				renderLogin(w, r, templates, resp)
				return
			}

			resp.Summary = "请求网络错误,请刷新重试"
			renderLogin(w, r, templates, resp)
			return
		}

		err = pwd.BcryptComparePassword(user.HashedPassword, password)
		if err != nil {
			resp.Summary = "邮箱或密码错误"
			renderLogin(w, r, templates, resp)
			return
		}

		encoded, err := secureCookie.Encode(AuthorizeCookie, &Authorize{ID: user.ID, Name: user.Username})
		if err != nil {
			resp.Summary = "请求网络错误,请刷新重试(cookie)"
			renderLogin(w, r, templates, resp)
			return
		}

		c := cookie.NewCookie(cookie.AuthorizeName, encoded, time.Now().Add(time.Duration(7200)*time.Second))
		http.SetCookie(w, c)
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func Logout(templates fs.FS) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie.DeleteCookie(w, cookie.AuthorizeName)
		http.Redirect(w, r, "/login", http.StatusFound)
	}
}

func renderRegister(w http.ResponseWriter, r *http.Request, templates fs.FS, data any) {
	renderLayout(w, r, templates, data, "user/register.html.tmpl")
}

func renderLogin(w http.ResponseWriter, r *http.Request, templates fs.FS, data any) {
	renderLayout(w, r, templates, data, "user/login.html.tmpl")
}

func viladatorRegister(email, username, password string) (registerPageData, bool) {
	ok := true
	resp := registerPageData{
		Email:    email,
		Username: username,
		Password: password,
	}

	if !ValidateRxEmail(email) {
		resp.EmailMsg = "请填写正确的邮箱地址"
		ok = false
	}
	if !ValidateRxUsername(username) {
		resp.UsernameMsg = "名称(6-20,字母,数字)"
		ok = false
	}
	if !ValidatePassword(password) {
		resp.PasswordMsg = "密码(8-20位)"
		ok = false
	}

	return resp, ok
}

func viladatorLogin(email, password string) (loginPageData, bool) {
	ok := true
	errs := loginPageData{
		Email:    email,
		Password: password,
	}

	if !ValidateRxEmail(email) {
		errs.EmailMsg = "请填写正确的邮箱地址"
		ok = false
	}
	if len(password) == 0 {
		errs.PasswordMsg = "请填写正确的密码"
		ok = false
	}

	return errs, ok
}
