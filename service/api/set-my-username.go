package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"cloudy-pics.uniroma1.it/cloudy-pics/service/api/reqcontext"
	"cloudy-pics.uniroma1.it/cloudy-pics/service/database_nosql"
)

func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var username Username
	err := json.NewDecoder(r.Body).Decode(&username)
	if err != nil {
		// The body was not a parseable JSON, reject it
		w.WriteHeader(http.StatusBadRequest)
		return
	} else if username.isNotValid() {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = rt.db.SetMyUsername(ps.ByName("user_id"), username.Name)
	if errors.Is(err, database_nosql.ErrUsernameAlreadyTaken) {
		w.WriteHeader(http.StatusConflict)
		return
	} else if err != nil {
		// some internal problem with the database
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
