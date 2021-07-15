package httpserver

import (
	"context"
	"errors"
	"html/template"
	"net/http"
	"path"
	"time"

	"github.com/gorilla/mux"
	"github.com/tyghr/social_network/internal/model"
)

var (
	userPageTmpl    = "http_tmpl/user_page.tmpl"
	userProfileTmpl = "http_tmpl/user_profile.tmpl"
)

type userPage struct {
	Title        string
	SelfUserName string
	UserName     string
	Posts        []model.Post
}

type userProfile struct {
	Title        string
	SelfUserName string
	UserName     string
	Name         string
	SurName      string
	Gender       string
	Age          int
	BirthDate    time.Time
	City         string
	Interests    []string
	Friends      []string
}

func (s *Server) showUserPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//w.Header().Set("Content-Type", "application/json; charset=utf8")

		ctx, cancel := context.WithTimeout(r.Context(), timeoutDefault)
		defer cancel()

		vars := mux.Vars(r)
		uP, ok := vars["user"]
		if !ok {
			s.error(w, r, http.StatusUnprocessableEntity, errors.New("path_user is wrong")) //422
			return
		}
		s.logger.Debugf("showUserPage query received (%s)", uP)

		selfUserName := ctx.Value(ctxKeyUserName).(string)

		posts, err := s.db.GetPosts(ctx, model.PostFilter{
			UserName: uP,
			// Header      string `json:"header"`
			// Text        string `json:"text"`
			// UpdatedFrom string `json:"updated_from"`
			// UpdatedTo   string `json:"updated_to"`
			// PageNum     int    `json:"pagenum"`
			// PageSize    int    `json:"pagesize"`
		})
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		t := template.Must(template.New(path.Base(userPageTmpl)).ParseFiles(userPageTmpl))
		err = t.Execute(w, userPage{
			Title:        globalTitle,
			SelfUserName: selfUserName,
			UserName:     uP,
			Posts:        posts,
		})
		if err != nil {
			s.logger.Errorf("failed render user page template: %v", err)
		}
	}
}

func (s *Server) showProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//w.Header().Set("Content-Type", "application/json; charset=utf8")

		ctx, cancel := context.WithTimeout(r.Context(), timeoutDefault)
		defer cancel()

		vars := mux.Vars(r)
		uP, ok := vars["user"]
		if !ok {
			s.error(w, r, http.StatusUnprocessableEntity, errors.New("path_user is wrong")) //422
			return
		}
		s.logger.Debugf("showProfile query received (%s)", uP)

		selfUserName := ctx.Value(ctxKeyUserName).(string)

		profile, err := s.db.GetProfile(ctx, uP)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		pGender := "M"
		if !profile.Gender {
			pGender = "F"
		}

		t := template.Must(template.New(path.Base(userProfileTmpl)).ParseFiles(userProfileTmpl))
		err = t.Execute(w, userProfile{
			Title:        globalTitle,
			SelfUserName: selfUserName,
			UserName:     profile.UserName,
			Name:         profile.Name,
			SurName:      profile.SurName,
			Gender:       pGender,
			Age:          int(time.Since(profile.BirthDate).Hours() / 24 / 365.25),
			BirthDate:    profile.BirthDate,
			City:         profile.City,
			Interests:    profile.Interests,
			Friends:      profile.Friends,
		})
		if err != nil {
			s.logger.Errorf("failed render user rofile template: %v", err)
		}
	}
}

func (s *Server) addFriend() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//w.Header().Set("Content-Type", "application/json; charset=utf8")

		ctx, cancel := context.WithTimeout(r.Context(), timeoutDefault)
		defer cancel()

		vars := mux.Vars(r)
		uP, ok := vars["user"]
		if !ok {
			s.error(w, r, http.StatusUnprocessableEntity, errors.New("path_user is wrong")) //422
			return
		}

		selfUserName := ctx.Value(ctxKeyUserName).(string)
		s.logger.Debugf("addFriend query received (%s want to add %s)", selfUserName, uP)

		err := s.db.AddFriend(ctx, selfUserName, uP)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		http.Redirect(w, r, "/user/"+uP+"/profile", http.StatusFound)
	}
}