package httpserver

import "net/http"

func (s *Server) showIndex() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ok, u := getUserSession(r)
		if !ok || u == "" {
			http.Redirect(w, r, "/login", http.StatusFound)
		} else {
			http.Redirect(w, r, "/user/"+u, http.StatusFound)
		}
	})
}
