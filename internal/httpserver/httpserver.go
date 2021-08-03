package httpserver

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/tyghr/logger"
	"github.com/tyghr/social_network/internal/config"
	"github.com/tyghr/social_network/internal/storage"

	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	// _ "github.com/tyghr/social_network/internal/api/docs"
)

const (
	ctxKeyRequestID ctxKey = iota
	ctxKeyUserName
)

var (
	timeoutDefault = 30 * time.Second
	globalTitle    = "Social network"
)

type ctxKey int8

type Server struct {
	db     storage.DataBase
	router *mux.Router
	conf   *config.Config
	logger logger.Logger
}

func NewServer(db storage.DataBase, conf *config.Config, l logger.Logger) *Server {
	s := &Server{
		router: mux.NewRouter(),
		db:     db,
		conf:   conf,
		logger: l,
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
	s.router.Use(s.setRequestID)
	s.router.Use(s.logRequest)
	s.router.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"})))

	attachProfiler(s.router)

	s.router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"), //The url pointing to API definition
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("#swagger-ui"),
	))

	s.router.HandleFunc("/", s.showIndex()).Methods(http.MethodGet)

	s.router.HandleFunc("/register", s.showRegister()).Methods(http.MethodGet)
	s.router.HandleFunc("/register", s.register()).Methods(http.MethodPost)
	s.router.HandleFunc("/search/user", s.searchUser()).Methods(http.MethodPost) // TMP

	s.router.HandleFunc("/logout", s.logout()).Methods(http.MethodPost)
	s.router.HandleFunc("/login", s.login()).Methods(http.MethodPost)
	s.router.HandleFunc("/login", s.showLogin()).Methods(http.MethodGet)

	s.router.HandleFunc("/post/edit", s.authSession(s.showEditPost())).Methods(http.MethodGet)
	s.router.HandleFunc("/post/edit", s.authSession(s.editPost())).Methods(http.MethodPost)
	s.router.HandleFunc("/post/delete", s.authSession(s.deletePost())).Methods(http.MethodPost)

	s.router.HandleFunc("/search/user", s.authSession(s.showSearchUser())).Methods(http.MethodGet)
	// s.router.HandleFunc("/search/user", s.authSession(s.searchUser())).Methods(http.MethodPost)

	s.router.HandleFunc("/user/{user}", s.authSession(s.showUserPage())).Methods(http.MethodGet)
	s.router.HandleFunc("/user/{user}/profile", s.authSession(s.showProfile())).Methods(http.MethodGet)
	s.router.HandleFunc("/user/{user}/add_to_friends", s.authSession(s.addFriend())).Methods(http.MethodPost)

	http.Handle("/", s.router)
}

func (s *Server) setRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyRequestID, id)))
	})
}

func (s *Server) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.logger.Infof("started %s %s remote_addr:%v request_id:%v", r.Method, r.RequestURI, r.RemoteAddr, r.Context().Value(ctxKeyRequestID))

		start := time.Now()
		rw := &responseWriter{w, http.StatusOK}
		next.ServeHTTP(rw, r)

		s.logger.Infof(
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
