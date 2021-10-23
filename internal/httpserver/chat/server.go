package chat

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/tyghr/logger"
	config "github.com/tyghr/social_network/internal/config/chat"
	"github.com/tyghr/social_network/internal/storage"
)

type Chat struct {
	store  storage.Chat
	router *mux.Router
	conf   *config.Config
	logger logger.Logger
}

func NewChatServer(storeChat storage.Chat, conf *config.Config, l logger.Logger) *Chat {
	chat := &Chat{
		router: mux.NewRouter(),
		store:  storeChat,
		conf:   conf,
		logger: l,
	}

	chat.router.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"})))

	chat.router.HandleFunc("/ws/chat", chat.handleChatWS)
	// chat.router.HandleFunc("/chat", chat.handleChatHome).Methods(http.MethodGet)

	http.Handle("/", chat.router)

	return chat
}

func (chat *Chat) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	chat.router.ServeHTTP(w, r)
}
