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

func (rt *_router) getPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	user_id := ps.ByName("photo_id")

	var photo Photo
	dbphoto, err := rt.db.GetPhoto(user_id, strings.Split(r.Header.Get("Authorization"), "Bearer ")[1])
	if errors.Is(err, database.ErrPhotoDoesNotExist) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if errors.Is(err, database.ErrBanned) {
		w.WriteHeader(http.StatusForbidden)
		return
	} else if err != nil {
		// some internal problem with the database
		ctx.Logger.WithError(err).Error("can't get the photo")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	photo.FromDatabase(dbphoto)

	w.Header().Set("content-type", "application/json")
	err = json.NewEncoder(w).Encode(photo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
