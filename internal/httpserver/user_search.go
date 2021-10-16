package httpserver

import (
	"context"
	"errors"
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
	Title         string
	SelfUserName  string
	UserName      string
	Name          string
	SurName       string
	Gender        string
	BirthDateFrom time.Time
	BirthDateTo   time.Time
	AgeFrom       string
	AgeTo         string
	City          string
	Interests     string
	Friends       string
	FoundUsers    []foundUser
}

type foundUser struct {
	UserName  string
	Name      string
	SurName   string
	Age       int
	BirthDate time.Time
	Gender    string
	City      string
}

type searchUser struct {
	UserName      string    `schema:"username"`
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

		ctx, cancel := context.WithTimeout(r.Context(), timeoutDefault)
		defer cancel()

		s.logger.Debugf("showSearchUser query received")

		selfUserName := ctx.Value(ctxKeyUserName).(string)

		t := template.Must(template.New(path.Base(userSearchTmpl)).ParseFiles(userSearchTmpl))
		err := t.Execute(w, searchUserPage{
			Title:        globalTitle,
			SelfUserName: selfUserName,
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

		selfUserName, ok := ctx.Value(ctxKeyUserName).(string)
		if !ok {
			s.error(w, r, http.StatusUnauthorized, errors.New("username is empty"))
			return
		}

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

		users, err := s.stor.DB().SearchUser(ctx, model.UserFilter{
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

		foundUsers := []foundUser{}
		for _, fu := range users {
			fuGender := "M"
			if !fu.Gender {
				fuGender = "F"
			}
			foundUsers = append(foundUsers, foundUser{
				UserName:  fu.UserName,
				Name:      fu.Name,
				SurName:   fu.SurName,
				Age:       int(time.Since(fu.BirthDate).Hours() / 24 / 365.25),
				BirthDate: fu.BirthDate,
				Gender:    fuGender,
				City:      fu.City,
			})
		}

		t := template.Must(template.New(path.Base(userSearchTmpl)).ParseFiles(userSearchTmpl))
		err = t.Execute(w, searchUserPage{
			Title:         globalTitle,
			SelfUserName:  selfUserName,
			UserName:      u.UserName,
			Name:          u.Name,
			SurName:       u.SurName,
			Gender:        u.Gender,
			BirthDateFrom: u.BirthDateFrom,
			BirthDateTo:   u.BirthDateTo,
			AgeFrom:       u.AgeFrom,
			AgeTo:         u.AgeTo,
			City:          u.City,
			Interests:     u.Interests,
			Friends:       u.Friends,
			FoundUsers:    foundUsers,
		})
		if err != nil {
			s.logger.Errorf("failed render user search template: %v", err)
		}
	}
}
