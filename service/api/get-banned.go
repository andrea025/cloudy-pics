package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"cloudy-pics.uniroma1.it/cloudy-pics/service/api/reqcontext"
	"cloudy-pics.uniroma1.it/cloudy-pics/service/database_nosql"
)

func (rt *_router) getBanned(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	ret := []UserShortInfo{}

	dbbanned, err := rt.db.GetBanned(ps.ByName("user_id"), strings.Split(r.Header.Get("Authorization"), "Bearer ")[1])
	if errors.Is(err, database_nosql.ErrUserDoesNotExist) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if errors.Is(err, database_nosql.ErrBanned) {
		w.WriteHeader(http.StatusForbidden)
		return
	} else if err != nil {
		// In this case, we have an error on our side. Log the error (so we can be notified) and send a 500 to the user
		// Note: we are using the "logger" inside the "ctx" (context) because the scope of this issue is the request.
		ctx.Logger.WithError(err).Error("can't list users this user has banned")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	for _, dbuser := range dbbanned {
		var user UserShortInfo
		user.FromDatabase(dbuser)
		ret = append(ret, user)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(UsersInfoListResponse{Data: ret})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
