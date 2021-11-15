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
	IsFriend     bool
	Subscribed   bool
}

type userProfile struct {
	Title         string
	SelfUserName  string
	UserName      string
	Name          string
	SurName       string
	Age           int
	BirthDate     time.Time
	Gender        string
	City          string
	Interests     []string
	Friends       []string
	Subscriptions []string
	Subscribers   []string
	IsFriend      bool
	Subscribed    bool
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

		posts, err := s.stor.DB().GetPosts(ctx, model.PostFilter{
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

		isF, isS, err := s.stor.DB().GetRelations(ctx, selfUserName, uP)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		t := template.Must(template.New(path.Base(userPageTmpl)).ParseFiles(userPageTmpl))
		err = t.Execute(w, userPage{
			Title:        globalTitle,
			SelfUserName: selfUserName,
			UserName:     uP,
			IsFriend:     isF,
			Subscribed:   isS,
			Posts:        posts,
		})
		if err != nil {
			s.logger.Errorf("failed render user page template: %v", err)
		}

		// reset cursor_counter (rabbit)
		err = s.stor.Q().UpdateCursorCounter(ctx, selfUserName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
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

		profile, err := s.stor.DB().GetProfile(ctx, uP)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		pGender := "M"
		if !profile.Gender {
			pGender = "F"
		}

		isF, isS, err := s.stor.DB().GetRelations(ctx, selfUserName, uP)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		t := template.Must(template.New(path.Base(userProfileTmpl)).ParseFiles(userProfileTmpl))
		err = t.Execute(w, userProfile{
			Title:         globalTitle,
			SelfUserName:  selfUserName,
			UserName:      profile.UserName,
			Name:          profile.Name,
			SurName:       profile.SurName,
			Gender:        pGender,
			Age:           int(time.Since(profile.BirthDate).Hours() / 24 / 365.25),
			BirthDate:     profile.BirthDate,
			City:          profile.City,
			Interests:     profile.Interests,
			Friends:       profile.Friends,
			Subscriptions: profile.Subscriptions,
			Subscribers:   profile.Subscribers,
			IsFriend:      isF,
			Subscribed:    isS,
		})
		if err != nil {
			s.logger.Errorf("failed render user profile template: %v", err)
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

		err := s.stor.DB().AddFriend(ctx, selfUserName, uP)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		http.Redirect(w, r, "/user/"+uP+"/profile", http.StatusFound)
	}
}

func (s *Server) subscribe() http.HandlerFunc {
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
		s.logger.Debugf("subscribe query received (%s subscribing to %s)", selfUserName, uP)

		err := s.stor.DB().Subscribe(ctx, selfUserName, uP)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		// add sub to rebuild queue (rabbit)
		err = s.stor.Q().PostRebuildSubsFeedRequest(ctx, []string{selfUserName})
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		// TODO change redirection
		http.Redirect(w, r, "/user/"+uP+"/profile", http.StatusFound)
	}
}

func (s *Server) unsubscribe() http.HandlerFunc {
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
		s.logger.Debugf("unsubscribe query received (%s unsubscribing from %s)", selfUserName, uP)

		err := s.stor.DB().Unsubscribe(ctx, selfUserName, uP)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		// add sub to rebuild queue (rabbit)
		err = s.stor.Q().PostRebuildSubsFeedRequest(ctx, []string{selfUserName})
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		// TODO change redirection
		http.Redirect(w, r, "/user/"+uP+"/profile", http.StatusFound)
	}
}
