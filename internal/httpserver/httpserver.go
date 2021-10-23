package httpserver

import (
	"encoding/json"
	"net/http"
	"net/http/pprof"
	"sync"
	"time"

	"github.com/tyghr/logger"
	"github.com/tyghr/social_network/internal/config"
	"github.com/tyghr/social_network/internal/storage"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	// _ "github.com/tyghr/social_network/internal/api/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

const (
	ctxKeyRequestID ctxKey = iota
	ctxKeyUserName
)

var (
	timeoutDefault       = 30 * time.Second
	globalTitle          = "Social network"
	wsSecureProtocolType = "SPTI"
)

type ctxKey int8

type Server struct {
	stor   *storage.Storage
	router *mux.Router
	conf   *config.Config
	logger logger.Logger

	reqIDs map[string]struct{}
	rlock  *sync.Mutex
}

func NewServer(stor *storage.Storage, conf *config.Config, l logger.Logger) *Server {
	s := &Server{
		router: mux.NewRouter(),
		stor:   stor,
		conf:   conf,
		logger: l,
		reqIDs: make(map[string]struct{}),
		rlock:  &sync.Mutex{},
	}
	s.configureRouter()

	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func attachProfiler(router *mux.Router) {
	router.HandleFunc("/debug/pprof/", pprof.Index)
	router.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	router.HandleFunc("/debug/pprof/profile", pprof.Profile)
	router.HandleFunc("/debug/pprof/symbol", pprof.Symbol)

	// Manually add support for paths linked to by index page at /debug/pprof/
	router.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	router.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	router.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	router.Handle("/debug/pprof/block", pprof.Handler("block"))
}

func (s *Server) configureRouter() {
	s.router.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"})))
	attachProfiler(s.router)

	swaggerRouter := s.router.PathPrefix("/swagger").Subrouter()
	swaggerRouter.Path("/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"), //The url pointing to API definition
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("#swagger-ui"),
	))

	wsRouter := s.router.PathPrefix("/ws").Subrouter()
	wsRouter.HandleFunc("/feed", s.authSession(s.handleFeedWS()))

	mainRouter := s.router.PathPrefix("/").Subrouter()
	mainRouter.Use(s.setRequestIDHandler)
	mainRouter.Use(s.logRequest)

	mainRouter.HandleFunc("/register", s.showRegister()).Methods(http.MethodGet)
	mainRouter.HandleFunc("/register", s.register()).Methods(http.MethodPost)

	mainRouter.HandleFunc("/logout", s.logout()).Methods(http.MethodPost)
	mainRouter.HandleFunc("/login", s.login()).Methods(http.MethodPost)
	mainRouter.HandleFunc("/login", s.showLogin()).Methods(http.MethodGet)

	mainRouter.HandleFunc("/post/edit", s.authSession(s.showUpsertPost())).Methods(http.MethodGet)
	mainRouter.HandleFunc("/post/edit", s.authSession(s.upsertPost())).Methods(http.MethodPost)
	mainRouter.HandleFunc("/post/delete", s.authSession(s.deletePost())).Methods(http.MethodPost)

	mainRouter.HandleFunc("/search/user", s.authSession(s.showSearchUser())).Methods(http.MethodGet)
	mainRouter.HandleFunc("/search/user", s.authSession(s.searchUser())).Methods(http.MethodPost)
	// mainRouter.HandleFunc("/search/user", s.searchUser()).Methods(http.MethodPost) // TMP

	mainRouter.HandleFunc("/user/{user}", s.authSession(s.showUserPage())).Methods(http.MethodGet)
	mainRouter.HandleFunc("/user/{user}/profile", s.authSession(s.showProfile())).Methods(http.MethodGet)
	mainRouter.HandleFunc("/user/{user}/add_to_friends", s.authSession(s.addFriend())).Methods(http.MethodPost)

	mainRouter.HandleFunc("/user/{user}/subscribe", s.authSession(s.subscribe())).Methods(http.MethodPost)
	mainRouter.HandleFunc("/user/{user}/unsubscribe", s.authSession(s.unsubscribe())).Methods(http.MethodPost)
	mainRouter.HandleFunc("/feed", s.authSession(s.showFeedPage())).Methods(http.MethodGet)

	mainRouter.HandleFunc("/chat/{user}", s.authSession(s.handleChatHome())).Methods(http.MethodGet)

	mainRouter.HandleFunc("/session_validate", s.checkRequestIDHandler()).Methods(http.MethodPost)

	mainRouter.HandleFunc("/", s.showIndex()).Methods(http.MethodGet)

	http.Handle("/", s.router)
}

func (s *Server) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.logger.Debugf("started %s %s remote_addr:%v request_id:%v", r.Method, r.RequestURI, r.RemoteAddr, r.Context().Value(ctxKeyRequestID))

		start := time.Now()
		rw := &responseWriter{w, http.StatusOK}
		next.ServeHTTP(rw, r)

		s.logger.Debugf(
			"completed with %d %s in %v. remote_addr:%v request_id:%v",
			rw.code,
			http.StatusText(rw.code),
			time.Since(start),
			r.RemoteAddr,
			r.Context().Value(ctxKeyRequestID),
		)
	})
}

func (s *Server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.logger.Error(err)
	if _, ok := err.(json.Marshaler); ok {
		s.respond(w, r, code, map[string]interface{}{"error": err})
		return
	}
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (s *Server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}
