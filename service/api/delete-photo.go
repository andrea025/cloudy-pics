package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"wasa-photo.uniroma1.it/wasa-photo/service/api/reqcontext"
	"wasa-photo.uniroma1.it/wasa-photo/service/database_nosql"
)

func (rt *_router) deletePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	pid := ps.ByName("photo_id")
	req_id := strings.Split(r.Header.Get("Authorization"), "Bearer ")[1]

	url, err := rt.db.DeletePhoto(pid, req_id)
	if errors.Is(err, database_nosql.ErrPhotoDoesNotExist) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if errors.Is(err, database_nosql.ErrDeletePhotoForbidden) {
		w.WriteHeader(http.StatusForbidden)
		return
	} else if err != nil {
		ctx.Logger.WithError(err).WithField("id", pid).Error("can't delete the photo form the database")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	
	// Delete the photo from s3 bucket
	err = rt.s3.DeletePhoto(url)

	if err != nil {
	    ctx.Logger.WithError(err).WithField("id", pid).Error("can't delete the photo form s3 bucket")
		w.WriteHeader(http.StatusInternalServerError)
	    return
	}

	w.WriteHeader(http.StatusNoContent)
}
