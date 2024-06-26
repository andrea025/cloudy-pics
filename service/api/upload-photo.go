package api

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"wasa-photo.uniroma1.it/wasa-photo/service/api/reqcontext"
	"wasa-photo.uniroma1.it/wasa-photo/service/database"
)

var storageBasePath string = "./storage/"

type Photo struct {
	Id              string             `json:"id"`
	CreatedDatetime string             `json:"created_datetime"`
	PhotoUrl        string             `json:"photo_url"`
	Owner           UserShortInfo      `json:"owner"`
	Likes           LikesCollection    `json:"likes"`
	Comments        CommentsCollection `json:"comments"`
}

func (p *Photo) FromDatabase(photo database.Photo) {
	p.Id = photo.Id
	p.CreatedDatetime = photo.CreatedDatetime
	p.PhotoUrl = photo.PhotoUrl
	p.Owner.FromDatabase(photo.Owner)
	p.Likes.FromDatabase(photo.Likes)
	p.Comments.FromDatabase(photo.Comments)
}

type UserShortInfo struct {
	Id       string `json:"id"`
	Username string `json:"username"`
}

func (usi *UserShortInfo) FromDatabase(userShortInfo database.UserShortInfo) {
	usi.Id = userShortInfo.Id
	usi.Username = userShortInfo.Username
}

type LikesCollection struct {
	Count int             `json:"count"`
	Users []UserShortInfo `json:"users"`
}

func (lc *LikesCollection) FromDatabase(lcdb database.LikesCollection) {
	lc.Count = lcdb.Count
	lc.Users = []UserShortInfo{}
	for _, usi := range lcdb.Users {
		var user UserShortInfo
		user.FromDatabase(usi)
		lc.Users = append(lc.Users, user)
	}
}

type CommentsCollection struct {
	Count    int       `json:"count"`
	Comments []Comment `json:"comments"`
}

func (cc *CommentsCollection) FromDatabase(ccdb database.CommentsCollection) {
	cc.Count = ccdb.Count
	cc.Comments = []Comment{}
	for _, c := range ccdb.Comments {
		var comment Comment
		comment.FromDatabase(c)
		cc.Comments = append(cc.Comments, comment)
	}
}

func (rt *_router) uploadPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// make sure the binary file in the body request is an image (actually it is not a proper check for its validity, since the client is taking the image/png from the extension of the file being uploaded, but fair enough)
	ctype := strings.Split(r.Header.Get("Content-type"), ";")[0]
	if ctype == "" || !(ctype == "image/png" || ctype == "image/jpeg") {
		// the request has no Content-type header, therefore is not valid
		// rt.baseLogger.Warning("uploadPhoto: a request has been sent without Content-type header or with Content-type header different than image/png and image/jpeg, the binary string sent in the request body is not valid")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	clength := r.ContentLength
	if clength == -1 {
		// no Content-length, it is not valid because we need to check the length of the image before uploading it
		// rt.baseLogger.Warning("uploadPhoto: a request has been sent without Content-length header")
		w.WriteHeader(http.StatusBadRequest)
		return
	} else if clength > 1000000 { // 1 MB
		// rt.baseLogger.Warning("uploadPhoto: a request has been sent with a photo too big in size")
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		return
	}

	// reading the image
	body, err := io.ReadAll(r.Body)
	if err != nil {
		ctx.Logger.WithError(err).Error("can't read the binary string in the body of the request")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// generating  the photo id
	rdm, er := rand.Int(rand.Reader, big.NewInt(1000))
	if er != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Concatenating the user id with the base 2 representation of the random number generatend, and md5 hashing the result to get the photo id
	pid := fmt.Sprintf("%x", md5.Sum([]byte(strings.Split(r.Header.Get("Authorization"), "Bearer ")[1]+rdm.Text(2))))
	filename := pid + "." + strings.Split(ctype, "/")[1]

	file, e := os.Create(storageBasePath + filename)
	if e != nil {
		ctx.Logger.WithError(e).Error("can't create the file for the photo")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer file.Close()

	_, err = file.Write(body)
	if err != nil {
		ctx.Logger.WithError(err).Error("can't write the photo on the file")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// insert the new record into a table
	var photo Photo
	creation := (time.Now()).Format(time.RFC3339)
	creation_datetime := creation[0:10] + " " + creation[11:19] // correct format
	url := "/storage/" + filename

	dbphoto, erro := rt.db.UploadPhoto(pid, creation_datetime, url, strings.Split(r.Header.Get("Authorization"), "Bearer ")[1])
	if erro != nil {
		ctx.Logger.WithError(erro).Error("can't upload photo")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	photo.FromDatabase(dbphoto)

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(photo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
