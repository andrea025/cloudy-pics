package utils

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"sync"
	"time"
)

// State represents the different possible states in the system
type State int

const (
	Start State = iota
	PatchUsername 
	SearchUser
	GetFollowing
	LikePhoto
	UnlikePhoto
	FollowUser
	UnfollowUser
	UploadPhoto
	DeletePhoto
	GetPhoto
	GetUserProfile
	GetMyStream
)

func (s State) String() string {
	return [...]string{
		"Start",
		"PatchUsername",		
		"SearchUser",
		"GetFollowing",
		"LikePhoto",
		"UnlikePhoto",
		"FollowUser",
		"UnfollowUser",
		"UploadPhoto",
		"DeletePhoto", 
		"GetPhoto",
		"GetUserProfile",
		"GetMyStream",
	}[s]
}

// Global arrays holding user information
var Usernames = []string{
	"Amelia",
	"George",
	"Isla",
	"Noah",
	"Oliver",
	"Olivia",
}

var UserIDs = []string{
	"93ea6597c3cbd06e93a46b9f5368732d",
	"578ad8e10dc4edb52ff2bd4ec9bc93a3",
	"8f1e344a1a612bfd7418be5368b9f2fe",
	"cfa36b7c75e18a9dc6e2a35d19a58ee7",
	"27090706d42a2525b9a07222f68dd3d4",
	"ba546f8d6d55634ce9106423ee4c5275",
}

var PhotoIDs = []string{
	"09ba15df691c58ae2d29af6e57c6f2c1",
	"10239ff3625f47b35884433fc23c1d7b",
	"21cb9f48e66ef13fde5a356deaa76abd",
	"fbbb5b52dc9ca46079a1c397f085ab33",
	"4a5092490e9c6deeb32e0f1a6957b3e2",
	"0ae70a4da85a872f2aa4e5a89e896475",
}

type UserModel struct {
	username       string
	userId         string
	lastUserId     string
	lastPhotoId      string
	currentState     State
	durationsMutex   sync.Mutex
	durations        []int
	failedRequests   int
}

type UserModelStats struct {
	ValidRequests  int
	FailedRequests int
	TimeMean       float64
	TimeStdDev     float64
}

const debugLog = true

func log(args ...any) {
	if debugLog {
		fmt.Println(args...)
	}
}

func NewUserModel() *UserModel {
	generatedUsername := fmt.Sprint("user0", rand.Intn(100))

	return &UserModel{
		currentState:   Start,
		username:     generatedUsername,
		userId:  	  "",
		lastUserId:   "",
		lastPhotoId:  "",
		durationsMutex: sync.Mutex{},
		durations:      []int{},
		failedRequests: 0,
	}
}

func (u *UserModel) start() State {
	log("Entering Start state...")

	if u.userId == "" {
		id, err := u.collectHTTPRequestWithBody("/session", "", "POST", u.username)
		if err != nil {
			log(err)
			// Necessary becaue if the user is not created we do not have the bearer token to do the other requests
			os.Exit(1)
		}

		u.userId = id
	}

	rand.Seed(time.Now().UnixNano())
	next := rand.Float32()

	switch {
		case next < 0.20:
			return PatchUsername
		case next < 0.40:
			return SearchUser
		case next < 0.60:
			return GetFollowing
		case next < 0.70:
			return LikePhoto
		case next < 0.80:
			return FollowUser
		case next < 0.85:
			return UploadPhoto
		case next < 0.90:
			return GetPhoto
		case next < 0.95:
			return GetUserProfile
		default:
			return GetMyStream
	}
}

func (u *UserModel) patchUsername() State {
	log("Entering PatchUsername state...")

	randomInt := rand.Intn(100)
	newUsername := fmt.Sprintf("%s_%d", u.username, randomInt)	
	err := u.collectHTTPRequest("/users/" + u.userId, u.userId, "PATCH", newUsername)

	if err != nil {
		log(err)
	}

	return Start
}


func (u *UserModel) searchUser() State {
	log("Entering SearchUser state...")

	user := Usernames[rand.Intn(len(Usernames))]
	err := u.collectHTTPRequest("/users?username=" + user, u.userId, "GET", "")

	if err != nil {
		log(err)
	}

	return Start
}


func (u *UserModel) getUserProfile() State {
	log("Entering GetUserProfile state...")

	id := UserIDs[rand.Intn(len(UserIDs))]
	err := u.collectHTTPRequest("/users/" + id, u.userId, "GET", "")

	if err != nil {
		log(err)
	}

	return Start
}

func (u *UserModel) getFollowing() State {
	log("Entering GetFollowing state...")

	id := UserIDs[rand.Intn(len(UserIDs))]
	err := u.collectHTTPRequest("/users/" + id + "/following", u.userId, "GET", "")

	if err != nil {
		log(err)
	}

	return Start
}

func (u *UserModel) getPhoto() State {
	log("Entering GetPhoto state...")

	id := PhotoIDs[rand.Intn(len(PhotoIDs))]
	err := u.collectHTTPRequest("/photos/" + id, u.userId, "GET", "")

	if err != nil {
		log(err)
	}

	return Start
}

func (u *UserModel) getMyStream() State {
	log("Entering GetMyStream state...")
	err := u.collectHTTPRequest("/users/" + u.userId + "/stream", u.userId, "GET", "")

	if err != nil {
		log(err)
	}

	return Start
}

func (u *UserModel) likePhoto() State {
	log("Entering LikePhoto state...")

	id := PhotoIDs[rand.Intn(len(PhotoIDs))]
	err := u.collectHTTPRequest("/photos/" + id + "/likes/" + u.userId, u.userId, "PUT", "")

	if err != nil {
		log(err)
		return Start
	}

	u.lastPhotoId = id

	return UnlikePhoto
}


func (u *UserModel) unlikePhoto() State {
	log("Entering UnlikePhoto state...")

	err := u.collectHTTPRequest("/photos/" + u.lastPhotoId + "/likes/" + u.userId, u.userId, "DELETE", "")

	if err != nil {
		log(err)
	}

	return Start
}


func (u *UserModel) followUser() State {
	log("Entering FollowUser state...")

	id := UserIDs[rand.Intn(len(UserIDs))]

	// followUser
	err := u.collectHTTPRequest("/users/" + u.userId + "/following/" + id, u.userId, "PUT", "")

	if err != nil {
		log(err)
		return Start
	}

	u.lastUserId = id

	return UnfollowUser
}

func (u *UserModel) unfollowUser() State {
	log("Entering UnfollowUser state...")

	// unfollowUser
	err := u.collectHTTPRequest("/users/" + u.userId + "/following/" + u.lastUserId, u.userId, "DELETE", "")

	if err != nil {
		log(err)
	}

	return Start
}


func (u *UserModel) uploadPhoto() State {
	log("Entering UploadPhoto state...")

	// List all files in ./test-assets
	assetsPhotos, err := os.ReadDir("./pictures")
	if err != nil {
		panic(err)
	}

	// Pick a random photo
	photoToUpload := assetsPhotos[rand.Intn(len(assetsPhotos))].Name()
	photoPath := fmt.Sprintf("./pictures/%s", photoToUpload)
	id, duration, err := timeHTTPPostFile("/users/" + u.userId + "/photos", photoPath, photoToUpload, u.userId)
	if err != nil {
		u.collectDuration(-1)
		log(err)
		return Start
	}

	u.collectDuration(duration)

	u.lastPhotoId = id

	return DeletePhoto
}

func (u *UserModel) deletePhoto() State {
	log("Entering DeletePhoto state...")

	// deletePhoto
	err := u.collectHTTPRequest("/photos/" + u.lastPhotoId, u.userId, "DELETE", "")
	if err != nil {
		log(err)
	}

	return Start
}


func (u *UserModel) Run() {
	for {
		var nextState State

		switch u.currentState {
			case Start:
				nextState = u.start()
			case PatchUsername:
				nextState = u.patchUsername()
			case SearchUser:
				nextState = u.searchUser()
			case GetPhoto:
				nextState = u.getPhoto()
			case GetUserProfile:
				nextState = u.getUserProfile()
			case GetMyStream:
				nextState = u.getMyStream()
			case FollowUser:
				nextState = u.followUser()
			case UnfollowUser:
				nextState = u.unfollowUser()
			case LikePhoto:
				nextState = u.likePhoto()
			case UnlikePhoto:
				nextState = u.unlikePhoto()
			case UploadPhoto:
				nextState = u.uploadPhoto()
			case DeletePhoto:
				nextState = u.deletePhoto()
		}

		var lastRequestTime int
		u.durationsMutex.Lock()
		if len(u.durations) > 0 {
			lastRequestTime = u.durations[len(u.durations)-1]
		} else {
			lastRequestTime = 0
		}
		u.durationsMutex.Unlock()

		if nextState != Start && lastRequestTime > -1 {
			time.Sleep(2 * time.Second)
		}
		u.currentState = nextState
	}
}


func (u *UserModel) collectHTTPRequestWithBody(url string, bearer string, method string, body string) (string, error) {
	body, duration, err := TimeHTTPRequestWithBody(url, bearer, method, body)
	if err != nil {
		u.collectDuration(-1)
	} else {
		u.collectDuration(duration)
	}
	return body, err
}

func (u *UserModel) collectHTTPRequest(url string, bearer string, method string, body string) error {
	duration, err := TimeHTTPRequest(url, bearer, method, body)
	if err != nil {
		u.collectDuration(-1)
	} else {
		u.collectDuration(duration)
	}
	return err
}

func (u *UserModel) collectDuration(duration int) {
	u.durationsMutex.Lock()
	defer u.durationsMutex.Unlock()

	log("Collecting duration", duration)

	if duration == -1 {
		u.failedRequests++
	} else {
		u.durations = append(u.durations, duration)
	}
}

func (u *UserModel) ResetStatistics() UserModelStats {
	u.durationsMutex.Lock()
	defer u.durationsMutex.Unlock()

	// Compute the mean
	var sum int
	for _, d := range u.durations {
		sum += d
	}
	mean := float64(sum) / float64(len(u.durations))

	// Compute the standard deviation
	var sumSquaredDiff float64
	for _, d := range u.durations {
		diff := float64(d) - mean
		sumSquaredDiff += diff * diff
	}
	stdDev := math.Sqrt(sumSquaredDiff / float64(len(u.durations)))

	stats := UserModelStats{
		ValidRequests:  len(u.durations),
		FailedRequests: u.failedRequests,
		TimeMean:       mean,
		TimeStdDev:     stdDev,
	}

	// Reset the accumulated durations
	u.durations = []int{}
	u.failedRequests = 0

	return stats
}
