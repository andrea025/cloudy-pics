package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"cloudy-pics.uniroma1.it/cloudy-pics/service/api/reqcontext"
)

type PhotoListResponse struct {
	Data []Photo `json:"data"`
}

func (rt *_router) getMyStream(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	stream := []Photo{}
	uid := strings.Split(r.Header.Get("Authorization"), "Bearer ")[1]

	dbstream, err := rt.db.GetMyStream(uid)
	if err != nil {
		ctx.Logger.WithError(err).Error("can't get the stream")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, dbphoto := range dbstream {
		var photo Photo
		photo.FromDatabase(dbphoto)
		stream = append(stream, photo)
	}

	w.Header().Set("content-type", "application/json")
	err = json.NewEncoder(w).Encode(PhotoListResponse{Data: stream})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
