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

type User struct {
	Id             string  `json:"id"`
	Username       string  `json:"username"`
	Followers      int     `json:"followers"`
	Following      int     `json:"following"`
	UploadedPhotos int     `json:"uploaded_photos"`
	Photos         []Photo `json:"photos"`
}

func (u *User) FromDatabase(user database.User) {
	u.Id = user.Id
	u.Username = user.Username
	u.Followers = user.Followers
	u.Following = user.Following
	u.UploadedPhotos = user.UploadedPhotos
	u.Photos = []Photo{}
	for _, p := range user.Photos {
		var photo Photo
		photo.FromDatabase(p)
		u.Photos = append(u.Photos, photo)
	}
}

func (rt *_router) getUserProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	user_id := ps.ByName("user_id")

	var user User
	dbuser, err := rt.db.GetUserProfile(user_id, strings.Split(r.Header.Get("Authorization"), "Bearer ")[1])
	if errors.Is(err, database.ErrUserDoesNotExist) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if errors.Is(err, database.ErrBanned) {
		w.WriteHeader(http.StatusForbidden)
		return
	} else if err != nil {
		// some internal problem with the database
		ctx.Logger.WithError(err).Error("can't get the user profile")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user.FromDatabase(dbuser)

	w.Header().Set("content-type", "application/json")
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
