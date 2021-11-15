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
	postTmpl = "http_tmpl/post_edit.tmpl"
)

type postPage struct {
	Title        string
	SelfUserName string
}

type post struct {
	UserName string `schema:"username,required"`
	Header   string `schema:"post_name"`
	Text     string `schema:"text"`
}

func (s *Server) showUpsertPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//w.Header().Set("Content-Type", "application/json; charset=utf8")

		ctx, cancel := context.WithTimeout(r.Context(), timeoutDefault)
		defer cancel()

		selfUserName := ctx.Value(ctxKeyUserName).(string)
		s.logger.Debugw("showUpsertPost query received",
			"user", selfUserName)

		t := template.Must(template.New(path.Base(postTmpl)).ParseFiles(postTmpl))
		err := t.Execute(w, postPage{
			Title:        globalTitle,
			SelfUserName: selfUserName,
		})
		if err != nil {
			s.logger.Errorf("failed render post edit template: %v", err)
		}
	}
}

func (s *Server) upsertPost() http.HandlerFunc {
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
		var p post
		err = decoder.Decode(&p, r.PostForm)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err) //400
			return
		}

		selfUserName := ctx.Value(ctxKeyUserName).(string)
		s.logger.Debugf("upsertPost query received (%s)", selfUserName)

		exist, err := s.stor.DB().UpsertPost(ctx, model.Post{
			UserName: selfUserName,
			Header:   p.Header,
			Text:     p.Text,
			Created:  time.Now(),
			Updated:  time.Now(),
		})
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		// get subscribers (mysql)
		subs, err := s.stor.DB().GetSubscribers(ctx, p.UserName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		subs = append(subs, selfUserName)

		s.logger.Debugw("upsertPost GetSubscribers",
			"user", selfUserName,
			"subs", subs)

		if !exist {
			// create tasks in queue (batch) (rabbit)
			err = s.stor.Q().AddPostBuckets(ctx, model.Post{
				UserName: p.UserName,
				Header:   p.Header,
				Text:     p.Text,
				Created:  time.Now(),
				Updated:  time.Now(),
			}, subs)
			if err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}

			// inc total_counters (rabbit)
			err = s.stor.Q().IncTotalCounters(ctx, subs)
			if err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}

		} else {
			// set rebuild_flag in DB
			err := s.stor.DB().SetFeedRebuildFlag(ctx, subs)
			if err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}

			// add subs to rebuild queue (rabbit)
			err = s.stor.Q().PostRebuildSubsFeedRequest(ctx, subs)
			if err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}
		}

		// .. processing queues

		http.Redirect(w, r, "/user/"+selfUserName, http.StatusFound)
	}
}

func (s *Server) deletePost() http.HandlerFunc {
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
		var p post
		err = decoder.Decode(&p, r.PostForm)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err) //400
			return
		}

		selfUserName := ctx.Value(ctxKeyUserName).(string)
		s.logger.Debugf("deletePost query received (%s)", selfUserName)

		err = s.stor.DB().DeletePost(ctx, model.Post{
			UserName: p.UserName,
			Header:   p.Header,
		})
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		// get subscribers (mysql)
		subs, err := s.stor.DB().GetSubscribers(ctx, p.UserName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		subs = append(subs, selfUserName)

		s.logger.Debugw("deletePost GetSubscribers",
			"user", selfUserName,
			"subs", subs)

		// inc cursor_counters (rabbit)
		err = s.stor.Q().IncCursorCounters(ctx, subs)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		// set rebuild_flag in DB
		err = s.stor.DB().SetFeedRebuildFlag(ctx, subs)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		// add subs to rebuild queue (rabbit)
		err = s.stor.Q().PostRebuildSubsFeedRequest(ctx, subs)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		http.Redirect(w, r, "/user/"+selfUserName, http.StatusFound)
	}
}
