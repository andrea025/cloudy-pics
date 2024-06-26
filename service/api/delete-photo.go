package api

import (
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/julienschmidt/httprouter"
	"wasa-photo.uniroma1.it/wasa-photo/service/api/reqcontext"
	"wasa-photo.uniroma1.it/wasa-photo/service/database"
)

func (rt *_router) deletePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	pid := ps.ByName("photo_id")
	req_id := strings.Split(r.Header.Get("Authorization"), "Bearer ")[1]

	url, err := rt.db.DeletePhoto(pid, req_id)
	if errors.Is(err, database.ErrPhotoDoesNotExist) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if errors.Is(err, database.ErrDeletePhotoForbidden) {
		w.WriteHeader(http.StatusForbidden)
		return
	} else if err != nil {
		ctx.Logger.WithError(err).WithField("id", pid).Error("can't delete the photo form the database")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	url_parts := strings.Split(url, "/")
	path := storageBasePath + url_parts[len(url_parts)-1]
	err = os.Remove(path)
	if err != nil {
		ctx.Logger.WithError(err).WithField("id", pid).Error("can't delete the photo form the filesystem")
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusNoContent)
}
