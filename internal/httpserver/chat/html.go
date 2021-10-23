package chat

import (
	"fmt"
	"html/template"
	"net/http"
	"path"
)

var (
	homeTemplateName     = "http_tmpl/chat.tmpl"
	homeTemplate         = template.Must(template.New(path.Base(homeTemplateName)).ParseFiles(homeTemplateName))
	chatEPTemplate       = "ws://%s/wschat?chat_id=%s"
	wsSecureProtocolType = "SPTI"
)

func (chat *Chat) handleChatHome(w http.ResponseWriter, r *http.Request) {
	chatID := r.URL.Query().Get("chat_id")
	if chatID == "" {
		http.Error(w, "empty chat_id", http.StatusBadRequest)
		return
	}

	err := homeTemplate.Execute(w, struct {
		Addr           string
		SecureProtocol string
	}{
		Addr:           fmt.Sprintf(chatEPTemplate, r.Host, chatID),
		SecureProtocol: wsSecureProtocolType,
	})
	if err != nil {
		panic(err)
	}
}
