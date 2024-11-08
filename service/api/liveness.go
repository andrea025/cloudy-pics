package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// liveness is an HTTP handler that checks the API server status. If the server cannot serve requests (e.g., some
// resources are not ready), this should reply with HTTP Status 500. Otherwise, with HTTP Status 200
func (rt *_router) liveness(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if err := rt.db.CheckConnectivity(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := rt.s3.CheckConnectivity(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := rt.lambda.CheckConnectivity(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	return
}
