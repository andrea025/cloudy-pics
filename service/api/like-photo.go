package api

import (
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"wasa-photo.uniroma1.it/wasa-photo/service/api/reqcontext"
	"wasa-photo.uniroma1.it/wasa-photo/service/database"
)

func (rt *_router) likePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	pid, uid := ps.ByName("photo_id"), ps.ByName("user_id")

	err := rt.db.LikePhoto(pid, uid)
	if errors.Is(err, database.ErrPhotoDoesNotExist) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if errors.Is(err, database.ErrBanned) || errors.Is(err, database.ErrNotSelfLike) {
		w.WriteHeader(http.StatusForbidden)
		return
	} else if err != nil {
		ctx.Logger.WithError(err).Error("can't create the like to the photo")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
