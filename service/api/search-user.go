package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"wasa-photo.uniroma1.it/wasa-photo/service/api/reqcontext"
	"wasa-photo.uniroma1.it/wasa-photo/service/database"
)

type UsersInfoListResponse struct {
	Data []UserShortInfo `json:"data"`
}

func (rt *_router) searchUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	ret := []UserShortInfo{}
	req_id := strings.Split(r.Header.Get("Authorization"), "Bearer ")[1]

	if r.URL.Query().Has("username") {
		var username Username
		username.Name = r.URL.Query().Get("username")
		if username.isNotValid() {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		dbuser, err := rt.db.SearchUser(username.Name, req_id)
		if errors.Is(err, database.ErrUserDoesNotExist) {
			w.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			// In this case, we have an error on our side. Log the error (so we can be notified) and send a 500 to the user
			// Note: we are using the "logger" inside the "ctx" (context) because the scope of this issue is the request.
			ctx.Logger.WithError(err).Error("can't search the user")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var user UserShortInfo
		user.FromDatabase(dbuser)
		ret = append(ret, user)
	} else {
		// Request an unfiltered list of fountains from the DB
		dbusers, err := rt.db.GetUsers(req_id)
		if err != nil {
			// In this case, we have an error on our side. Log the error (so we can be notified) and send a 500 to the user
			// Note: we are using the "logger" inside the "ctx" (context) because the scope of this issue is the request.
			ctx.Logger.WithError(err).Error("can't list users")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		for _, dbuser := range dbusers {
			var user UserShortInfo
			user.FromDatabase(dbuser)
			ret = append(ret, user)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(UsersInfoListResponse{Data: ret})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
