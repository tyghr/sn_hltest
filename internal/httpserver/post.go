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
	Title    string
	UserName string
}

type post struct {
	UserName string `schema:"username,required"`
	Header   string `schema:"post_name"`
	Text     string `schema:"text"`
}

func (s *Server) showEditPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//w.Header().Set("Content-Type", "application/json; charset=utf8")

		// ctx, cancel := context.WithTimeout(r.Context(), timeoutDefault)
		// defer cancel()

		s.logger.Debugf("showEditPost query received")

		t := template.Must(template.New(path.Base(postTmpl)).ParseFiles(postTmpl))
		err := t.Execute(w, postPage{
			Title:    globalTitle,
			UserName: "TODO",
		})
		if err != nil {
			s.logger.Errorf("failed render post edit template: %v", err)
		}
	}
}

func (s *Server) editPost() http.HandlerFunc {
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

		s.logger.Debugf("editPost query received")

		err = s.db.EditPost(ctx, model.Post{
			UserName: p.UserName,
			Header:   p.Header,
			Text:     p.Text,
			Updated:  time.Now(),
		})
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		// TODO show user page
		s.respond(w, r, http.StatusOK, "ok")
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

		s.logger.Debugf("deletePost query received")

		err = s.db.DeletePost(ctx, model.Post{
			UserName: p.UserName,
			Header:   p.Header,
		})
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		// TODO show user page
		s.respond(w, r, http.StatusOK, "ok")
	}
}
