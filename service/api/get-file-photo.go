package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"wasa-photo.uniroma1.it/wasa-photo/service/api/reqcontext"
)

func (rt *_router) getFilePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	http.ServeFile(w, r, storageBasePath+ps.ByName("filename"))
}
