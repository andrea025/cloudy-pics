package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"wasa-photo.uniroma1.it/wasa-photo/service/api/reqcontext"
	"wasa-photo.uniroma1.it/wasa-photo/service/database_nosql"
)

func (rt *_router) uncommentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	pid, cid, req_id := ps.ByName("photo_id"), ps.ByName("comment_id"), strings.Split(r.Header.Get("Authorization"), "Bearer ")[1]

	err := rt.db.UncommentPhoto(pid, cid, req_id)
	if errors.Is(err, database_nosql.ErrPhotoDoesNotExist) || errors.Is(err, database_nosql.ErrCommentDoesNotExist) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if errors.Is(err, database_nosql.ErrForbidden) {
		w.WriteHeader(http.StatusForbidden)
		return
	} else if err != nil {
		ctx.Logger.WithError(err).Error("can't delete the comment")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
