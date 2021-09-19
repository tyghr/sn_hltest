package chat

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/tyghr/social_network/internal/storage"
)

type Chat struct {
	store  storage.Chat
	router *mux.Router
}

func Init(storeChat storage.Chat) *Chat {
	chat := &Chat{
		router: mux.NewRouter(),
		store:  storeChat,
	}

	chat.router.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"})))

	chat.router.HandleFunc("/wschat", chat.handleChatWS)
	chat.router.HandleFunc("/chat", chat.handleChatHome).Methods(http.MethodGet)

	http.Handle("/", chat.router)

	return chat
}

func (chat *Chat) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	chat.router.ServeHTTP(w, r)
}
