package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"wasa-photo.uniroma1.it/wasa-photo/service/api/reqcontext"
)

type Username struct {
	Name string `json:"username"`
}

func (u *Username) isNotValid() bool {
	return len(u.Name) < 3 || len(u.Name) > 16
}

type UserId struct {
	Id string `json:"id"`
}

func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// query the User table to get a single row (id) where the username is equal to the one present in the body
	// if the query has that element, put it into the body in JSON format and return it
	// otherwise generate an id and perform an insert of the new user on the User table, and return the id like in the previous step

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

	var id string
	id, err = rt.db.DoLogin(username.Name)
	if err != nil {
		// some internal problem with the database
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("content-type", "application/json")
	err = json.NewEncoder(w).Encode(UserId{Id: id})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
