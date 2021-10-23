package httpserver

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"path"
	"time"

	"github.com/gorilla/websocket"
)

var (
	feedPageTmpl       = "http_tmpl/feed_page.tmpl"
	wsFeedAddrTemplate = "ws://%s/ws/feed"
)

type feedPage struct {
	Addr           string
	SecureProtocol string
	Title          string
	SelfUserName   string
}

func (s *Server) showFeedPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//w.Header().Set("Content-Type", "application/json; charset=utf8")
		// ctx, cancel := context.WithTimeout(r.Context(), timeoutDefault)
		// defer cancel()
		ctx := r.Context()

		ok, u := getUserSession(r)
		if !ok || u == "" {
			s.logger.Errorw("getUserSession", "error", errors.New("not authenticated"))
			return
		}

		selfUserName := ctx.Value(ctxKeyUserName).(string)
		s.logger.Debugf("showFeedPage query received (%s)", selfUserName)

		t := template.Must(template.New(path.Base(feedPageTmpl)).ParseFiles(feedPageTmpl))
		err := t.Execute(w, feedPage{
			Addr:           fmt.Sprintf(wsFeedAddrTemplate, r.Host),
			SecureProtocol: wsSecureProtocolType,
			Title:          globalTitle,
			SelfUserName:   selfUserName,
		})
		if err != nil {
			s.logger.Errorf("failed render feed page template: %v", err)
		}
	}
}

func (s *Server) handleFeedWS() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), timeoutDefault)
		defer cancel()

		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			s.logger.Errorw("upgrade", "error", err)
			return
		}
		defer c.Close()

		selfUserName := ctx.Value(ctxKeyUserName).(string)

		userToken, err := s.getUserToken(c, websocket.Subprotocols(r))
		if err != nil {
			s.logger.Errorw("getUserToken", "error", err)
		}

		// debug info
		s.logger.Debugw("handleFeedWS channel_info",
			"user", selfUserName,
			"userToken", userToken,
		)

		closeFunc := s.readingMessages(c, userToken)
		defer func() {
			closeFunc()
		}()

		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				s.logger.Debugw("ws closed by error:", "error", err)
				break
			}
			s.logger.Debugw("ws closed by msg", "message", string(message))
		}

		s.logger.Debugw("handleFeedWS stopped",
			"user", selfUserName,
			"userToken", userToken,
		)
	}
}

func (s *Server) readingMessages(c *websocket.Conn, userToken string) func() {
	chClose := make(chan struct{})
	close := func() {
		close(chClose)
	}

	go func() {
		id := "0"
		defer func() {
			s.logger.Debugw("stop watching redis",
				"userToken", userToken)
		}()
	L:
		for {
			select {
			case <-chClose:
				break L
			default:
			}

			s.logger.Debugw("invoking GetSubscriptionPosts")

			msgs, err := s.stor.C().GetSubscriptionPosts(context.TODO(), userToken, id)
			if err != nil {
				// TODO
				s.logger.Fatalw("GetSubscriptionPosts",
					"error", err)
			}
			for _, msg := range msgs {
				reset, ok := msg.Value("reset").(int)
				if ok && reset == 1 {
					s.logger.Debugw("reset FE msg")

					err = c.WriteMessage(websocket.TextMessage, []byte(`{"reset":1}`))
					if err != nil {
						s.logger.Errorw("write reset msg",
							"error", err)
						break
					}

					id = "0"
					continue
				}

				if msg.ID() > id {
					id = msg.ID()
				}

				pd, ok := msg.Value("post_data").(string)
				if !ok {
					s.logger.Errorw("no post_data in msg")
					break
				}

				s.logger.Debugw("read from redis")

				err = c.WriteMessage(websocket.TextMessage, []byte(pd))
				if err != nil {
					s.logger.Errorw("write msg",
						"error", err)
					break
				}
				s.logger.Debugw("send to ws")
			}

			time.Sleep(5 * time.Second)
		}
	}()

	// <close callBack>
	return close
}
