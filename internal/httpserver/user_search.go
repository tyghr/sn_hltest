package httpserver

import (
	"context"
	"html/template"
	"net/http"
	"path"
	"time"

	"github.com/gorilla/schema"
	"github.com/tyghr/social_network/internal/model"
)

var (
	userSearchTmpl = "http_tmpl/user_search.tmpl"
)

type searchUserPage struct {
	Title string
}

type searchUser struct {
	UserName      string    `schema:"username,required"`
	Name          string    `schema:"first_name"`
	SurName       string    `schema:"second_name"`
	Gender        string    `schema:"gender"` // ALL|M|F
	BirthDateFrom time.Time `schema:"birthdate_from"`
	BirthDateTo   time.Time `schema:"birthdate_to"`
	AgeFrom       string    `schema:"age_from"`
	AgeTo         string    `schema:"age_to"`
	City          string    `schema:"city"`
	Interests     string    `schema:"interests"`
	Friends       string    `schema:"friends"`
}

func (s *Server) showSearchUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//w.Header().Set("Content-Type", "application/json; charset=utf8")

		// ctx, cancel := context.WithTimeout(r.Context(), timeoutDefault)
		// defer cancel()

		s.logger.Debugf("showSearchUser query received")

		t := template.Must(template.New(path.Base(userSearchTmpl)).ParseFiles(userSearchTmpl))
		err := t.Execute(w, searchUserPage{
			Title: globalTitle,
		})
		if err != nil {
			s.logger.Errorf("failed render user search template: %v", err)
		}
	}
}

func (s *Server) searchUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//w.Header().Set("Content-Type", "application/json; charset=utf8")

		ctx, cancel := context.WithTimeout(r.Context(), timeoutDefault)
		defer cancel()

		err := r.ParseForm()
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err) //400
			return
		}

		decoder := schema.NewDecoder()
		var u searchUser
		err = decoder.Decode(&u, r.PostForm)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err) //400
			return
		}

		s.logger.Debugf("searchUser query received")

		users, err := s.db.SearchUser(ctx, model.UserFilter{
			UserName:      u.UserName,
			Name:          u.Name,
			SurName:       u.SurName,
			BirthDateFrom: u.BirthDateFrom,
			BirthDateTo:   u.BirthDateTo,
			AgeFrom:       u.AgeFrom,
			AgeTo:         u.AgeTo,
			Gender:        u.Gender,
			City:          u.City,
			Interests:     u.Interests,
			Friends:       u.Friends,
			//PageNum
			//PageSize
		})
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		// TODO show user list page
		s.respond(w, r, http.StatusOK, users)
	}
}
