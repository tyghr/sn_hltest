package counters

import (
	"context"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
)

func (srv *Srv) handleGetCounter() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//w.Header().Set("Content-Type", "application/json; charset=utf8")

		ctx, cancel := context.WithTimeout(r.Context(), timeoutDefault)
		defer cancel()

		vars := mux.Vars(r)
		uP, ok := vars["user"]
		if !ok {
			srv.error(w, r, http.StatusUnprocessableEntity, errors.New("path_user is wrong")) //422
			return
		}

		selfUserName := ctx.Value(ctxKeyUserName).(string)
		srv.logger.Debugf("handleGetCounter query received (%s subscribing to %s)", selfUserName, uP)

		count, err := srv.store.C().GetUnreadCount(ctx, selfUserName)
		if err != nil {
			srv.error(w, r, http.StatusInternalServerError, err)
			return
		}

		srv.respond(w, r, http.StatusOK, count)
	}
}
