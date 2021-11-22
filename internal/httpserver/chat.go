package httpserver

import (
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/gorilla/mux"
)

var (
	homeTemplateName = "chat.tmpl"
	chatEPTemplate   = "%s?chat_id=%s"
)

func (s *Server) handleChatHome() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		ok, u := getUserSession(r)
		if !ok || u == "" {
			s.logger.Errorw("getUserSession", "error", errors.New("not authenticated"))
			return
		}

		vars := mux.Vars(r)
		uP, ok := vars["user"]
		if !ok {
			s.error(w, r, http.StatusUnprocessableEntity, errors.New("path_user is wrong")) //422
			return
		}

		selfUserName := ctx.Value(ctxKeyUserName).(string)
		requestID := ctx.Value(ctxKeyRequestID).(string)
		s.logger.Debugw("handleChatHome query received",
			"username", selfUserName,
			"request_id", requestID,
			"chat_person", uP,
		)

		homeTemplate := s.getHtmlTemplate(homeTemplateName)

		err := homeTemplate.Execute(w, struct {
			Addr           string
			SecureProtocol string
			Title          string
			SelfUserName   string
			UserToken      string
		}{
			Addr:           fmt.Sprintf(chatEPTemplate, s.conf.ChatUrl, s.getChatID(selfUserName, uP)),
			SecureProtocol: wsSecureProtocolType,
			Title:          globalTitle,
			SelfUserName:   selfUserName,
			UserToken:      getToken(selfUserName, requestID),
		})
		if err != nil {
			s.logger.Errorf("failed render chat page template: %v", err)
		}
	}
}

func (srv *Server) getChatID(owner, person string) string {
	s := []string{owner, person}
	sort.Strings(s)

	hasher := sha512.New()
	hasher.Write([]byte(strings.Join(s, "_")))
	resp := hasher.Sum(nil)

	return base64.URLEncoding.EncodeToString(resp)
}

func getToken(selfUserName, requestID string) string {
	userToken := struct {
		Username  string `json:"username"`
		RequestID string `json:"request_id"`
	}{
		Username:  selfUserName,
		RequestID: requestID,
	}

	b, _ := json.Marshal(userToken)

	return base64.URLEncoding.EncodeToString(b)
}
