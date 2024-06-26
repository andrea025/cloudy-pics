package api

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"wasa-photo.uniroma1.it/wasa-photo/service/api/reqcontext"
	"wasa-photo.uniroma1.it/wasa-photo/service/database"
)

type Comment struct {
	Id              string        `json:"id"`
	Photo           string        `json:"photo"`
	User            UserShortInfo `json:"user"`
	Text            string        `json:"text"`
	CreatedDatetime string        `json:"created_datetime"`
}

func (c *Comment) FromDatabase(comment database.Comment) {
	c.Id = comment.Id
	c.Photo = comment.Photo
	c.User.FromDatabase(comment.User)
	c.Text = comment.Text
	c.CreatedDatetime = comment.CreatedDatetime
}

type CommentText struct {
	Text string `json:"text"`
}

func (ct *CommentText) isNotValid() bool {
	return len(ct.Text) < 3 || len(ct.Text) > 300
}

func (rt *_router) commentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var ct CommentText
	err := json.NewDecoder(r.Body).Decode(&ct)
	if err != nil {
		// The body was not a parseable JSON, reject it
		w.WriteHeader(http.StatusBadRequest)
		return
	} else if ct.isNotValid() {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var comment Comment
	pid := ps.ByName("photo_id")
	uid := strings.Split(r.Header.Get("Authorization"), "Bearer ")[1]
	creation := (time.Now()).Format(time.RFC3339)
	created_datetime := creation[0:10] + " " + creation[11:19] // correct format
	cid := fmt.Sprintf("%x", md5.Sum([]byte(pid+uid+creation)))
	// comment = Comment{Id: cid, Photo: pid, User: uid, Text: ct.Text, CreatedDatetime: creation_datetime}

	dbComment, er := rt.db.CommentPhoto(cid, pid, uid, ct.Text, created_datetime)
	if errors.Is(er, database.ErrBanned) {
		w.WriteHeader(http.StatusForbidden)
		return
	} else if errors.Is(er, database.ErrPhotoDoesNotExist) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if er != nil {
		// some internal problem with the database
		ctx.Logger.WithError(er).Error("can't create comment")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	comment.FromDatabase(dbComment)

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("content-type", "application/json")
	err = json.NewEncoder(w).Encode(comment)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
