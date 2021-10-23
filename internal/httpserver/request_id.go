package httpserver

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

func (s *Server) setRequestIDHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var keepReqID bool
		if strings.HasPrefix(r.URL.Path, "/chat/") {
			keepReqID = true
		}

		id := uuid.New().String()

		if keepReqID {
			s.putRequestID(id)
		}

		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyRequestID, id)))

		if keepReqID {
			go func() {
				time.Sleep(time.Second * 10)
				s.dropRequestID(id)
			}()
		}
	})
}

func (s *Server) putRequestID(reqID string) {
	s.rlock.Lock()
	defer s.rlock.Unlock()

	s.logger.Debugw("putRequestID",
		"req_id", reqID,
	)

	s.reqIDs[reqID] = struct{}{}
}

func (s *Server) dropRequestID(reqID string) {
	s.rlock.Lock()
	defer s.rlock.Unlock()

	s.logger.Debugw("dropRequestID",
		"req_id", reqID,
	)

	delete(s.reqIDs, reqID)
}

func (s *Server) popRequestID(reqID string) bool {
	s.rlock.Lock()
	defer s.rlock.Unlock()

	_, ok := s.reqIDs[reqID]

	s.logger.Debugw("popRequestID",
		"req_id", reqID,
		"ok", ok,
	)

	delete(s.reqIDs, reqID)

	return ok
}

func (s *Server) checkRequestIDHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf8")

		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err) //400
			return
		}

		j, err := base64.URLEncoding.DecodeString(string(b))
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err) //400
			return
		}

		var userToken struct {
			Username  string `json:"username"`
			RequestID string `json:"request_id"`
		}
		err = json.Unmarshal(j, &userToken)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err) //400
			return
		}

		s.respond(w, r, http.StatusOK, s.popRequestID(userToken.RequestID))
	}
}
