package counters

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/tyghr/logger"
	config "github.com/tyghr/social_network/internal/config/counters"
	"github.com/tyghr/social_network/internal/storage"
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

type Srv struct {
	store  *storage.CountersStorage
	router *mux.Router
	conf   *config.Config
	logger logger.Logger
}

func NewCountersServer(stor *storage.CountersStorage, conf *config.Config, l logger.Logger) *Srv {
	srv := &Srv{
		router: mux.NewRouter(),
		store:  stor,
		conf:   conf,
		logger: l,
	}

	srv.router.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"})))

	srv.router.Use(srv.setRequestIDHandler)
	srv.router.Use(srv.logRequest)

	srv.router.HandleFunc("/counter/{user}", srv.handleGetCounter()).Methods(http.MethodGet)

	srv.router.HandleFunc("/health_check", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods(http.MethodGet)

	http.Handle("/", srv.router)

	return srv
}

func (srv *Srv) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	srv.router.ServeHTTP(w, r)
}

func (srv *Srv) setRequestIDHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()

		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyRequestID, id)))
	})
}

func (srv *Srv) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		srv.logger.Debugf("started %s %s remote_addr:%v request_id:%v", r.Method, r.RequestURI, r.RemoteAddr, r.Context().Value(ctxKeyRequestID))

		start := time.Now()
		rw := &responseWriter{w, http.StatusOK}
		next.ServeHTTP(rw, r)

		srv.logger.Debugf(
			"completed with %d %s in %v. remote_addr:%v request_id:%v",
			rw.code,
			http.StatusText(rw.code),
			time.Since(start),
			r.RemoteAddr,
			r.Context().Value(ctxKeyRequestID),
		)
	})
}

func (srv *Srv) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	srv.logger.Error(err)
	if _, ok := err.(json.Marshaler); ok {
		srv.respond(w, r, code, map[string]interface{}{"error": err})
		return
	}
	srv.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (srv *Srv) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}
