package httpserver

import (
	"context"
	"crypto/sha256"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/gorilla/schema"
	"github.com/tyghr/social_network/internal/model"
)

var (
	registerTmpl = "register.tmpl"
)

type registerPage struct {
	Title string
}

type regUser struct {
	UserName  string    `schema:"username,required"`
	Name      string    `schema:"first_name"`
	SurName   string    `schema:"second_name"`
	Gender    string    `schema:"gender"` // M|F
	BirthDate time.Time `schema:"birthdate"`
	City      string    `schema:"city"`
	Interests string    `schema:"interests"`
	Password1 string    `schema:"password,required"`
	Password2 string    `schema:"repeat_password,required"`
}

var dateConverter = func(value string) reflect.Value {
	if v, err := time.Parse("2006-01-02", value); err == nil {
		return reflect.ValueOf(v)
	}
	return reflect.Value{}
}

func (s *Server) showRegister() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//w.Header().Set("Content-Type", "application/json; charset=utf8")

		// ctx, cancel := context.WithTimeout(r.Context(), timeoutDefault)
		// defer cancel()

		s.logger.Debugf("showRegister query received")

		t := s.getHtmlTemplate(registerTmpl)
		err := t.Execute(w, registerPage{
			Title: globalTitle,
		})
		if err != nil {
			s.logger.Errorf("failed render showRegister template: %v", err)
		}
	}
}

func (s *Server) register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//w.Header().Set("Content-Type", "application/json; charset=utf8")

		ctx, cancel := context.WithTimeout(r.Context(), timeoutDefault)
		defer cancel()

		err := r.ParseForm()
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err) //400
			return
		}
		defer r.Body.Close()

		decoder := schema.NewDecoder()
		decoder.RegisterConverter(time.Time{}, dateConverter)
		var u regUser
		err = decoder.Decode(&u, r.PostForm)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err) //400
			return
		}

		if u.Password1 != u.Password2 {
			s.error(w, r, http.StatusUnprocessableEntity, err) //422
			return
		}

		s.logger.Debugf("register query received")

		pHash := sha256.Sum256([]byte(u.Password1))

		interests := []string{}
		for _, i := range strings.Split(u.Interests, ",") {
			if ti := strings.TrimSpace(i); ti != "" {
				interests = append(interests, ti)
			}
		}

		err = s.stor.DB().Register(ctx, model.User{
			UserName:     u.UserName,
			PasswordHash: pHash[:],
			Name:         u.Name,
			SurName:      u.SurName,
			BirthDate:    u.BirthDate,
			Gender:       u.Gender == "M",
			City:         u.City,
			Interests:    interests,
		})
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		http.Redirect(w, r, "/user/"+u.UserName, http.StatusFound)
	}
}
