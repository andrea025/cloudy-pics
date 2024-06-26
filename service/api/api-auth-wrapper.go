package api

import (
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

// Checks if the user is authenticated (Bearer Authentication, ID in the Authorization header of the HTTP request)
// and allows the user to perform the requested action if no authorization is needed, otherwise the user shall also be authorized
// (the ID in the authorization header needs to be equal to the user_id path parameter) in order to proceed
func (rt *_router) auth(fn httprouter.Handle, authorization bool) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
		if len(authHeader) != 2 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		exists, err := rt.db.CheckUser(authHeader[1])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		} else if !exists {
			w.WriteHeader(http.StatusUnauthorized)
			return
		} else if auth_id, param := authHeader[1], ps.ByName("user_id"); (authorization && auth_id == param) || !authorization {
			// Delegate request to the given handle
			fn(w, r, ps)
		} else {
			w.WriteHeader(http.StatusForbidden)
		}
	}
}
