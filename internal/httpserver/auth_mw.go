package httpserver

import (
	"net/http"

	"crypto/sha256"
	"time"

	"github.com/gorilla/schema"
	"github.com/gorilla/securecookie"

	"context"

	"github.com/google/uuid"
)

var (
	cookieHandler = securecookie.New(
		securecookie.GenerateRandomKey(64),
		securecookie.GenerateRandomKey(32))
	cookieUserName   = "session_user"
	cookieUserExpire = 10 * time.Minute
	// realm               = "Please enter your credentials"
	// errNotAuthenticated = errors.New("not authenticated")
	// errCookieExpired    = errors.New("session expired")
	loginTmpl = "login.tmpl"
)

type loginPage struct {
	Title        string
	LoginLink    string
	RegisterLink string
}

type loginUser struct {
	UserName string `schema:"username,required"`
	Password string `schema:"password,required"`
}

func getUserSession(r *http.Request) (bool, string) {
	c, err := r.Cookie(cookieUserName)
	if err != nil {
		return false, ""
	}
	cookieValue := make(map[string]string)
	err = cookieHandler.Decode(cookieUserName, c.Value, &cookieValue)
	return err == nil, cookieValue["username"]
}

func (s *Server) setSession(w http.ResponseWriter, userName string) {
	exp := time.Now().Add(cookieUserExpire)
	sessionToken := uuid.New().String()
	value := map[string]string{
		"username": userName,
		"token":    sessionToken,
	}

	s.logger.Debugf("renew cookie for %s: %s (%s)", userName, sessionToken, exp.Format(time.RFC3339))

	if encoded, err := cookieHandler.Encode(cookieUserName, value); err == nil {
		http.SetCookie(w, &http.Cookie{
			Name:    cookieUserName,
			Value:   encoded,
			Path:    "/",
			Expires: exp,
		})
	}
}

func (s *Server) clearSession(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:   cookieUserName,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
}

func (s *Server) showLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//w.Header().Set("Content-Type", "application/json; charset=utf8")

		// ctx, cancel := context.WithTimeout(r.Context(), timeoutDefault)
		// defer cancel()

		s.logger.Debugf("showLogin query received")

		t := s.getHtmlTemplate(loginTmpl)
		err := t.Execute(w, loginPage{
			Title:        globalTitle,
			LoginLink:    "/login",
			RegisterLink: "/register",
		})
		if err != nil {
			s.logger.Errorf("failed render showLogin template: %v", err)
		}
	}
}

func (s *Server) login() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), timeoutDefault)
		defer cancel()

		err := r.ParseForm()
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err) //400
			return
		}

		decoder := schema.NewDecoder()
		var u loginUser
		err = decoder.Decode(&u, r.PostForm)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err) //400
			return
		}

		if u.UserName == "" || u.Password == "" {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		pHash := sha256.Sum256([]byte(u.Password))

		authOk, authErr := s.stor.DB().CheckAuth(ctx, u.UserName, pHash[:])
		if authErr != nil || !authOk {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		s.setSession(w, u.UserName)

		s.logger.Debugf("auth passed. user:%s", u.UserName)

		http.Redirect(w, r, "/user/"+u.UserName, http.StatusFound)
	})
}

func (s *Server) authSession(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ok, u := getUserSession(r)
		if !ok || u == "" {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		ctx := context.WithValue(r.Context(), ctxKeyUserName, u) //ctx.Value(ctxKeyUserName).(string)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *Server) logout() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.clearSession(w)
		http.Redirect(w, r, "/login", http.StatusFound)
	})
}
