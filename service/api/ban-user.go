package api

import (
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"cloudy-pics.uniroma1.it/cloudy-pics/service/api/reqcontext"
	"cloudy-pics.uniroma1.it/cloudy-pics/service/database_nosql"
)

func (rt *_router) banUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	uid, tuid := ps.ByName("user_id"), ps.ByName("target_user_id")

	err := rt.db.BanUser(uid, tuid)

	if errors.Is(err, database_nosql.ErrUserDoesNotExist) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		ctx.Logger.WithError(err).Error("can't create the banned relationship")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
