package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.POST("/session", rt.wrap(rt.doLogin))
	rt.router.GET("/photos/:photo_id", rt.auth(rt.wrap(rt.getPhoto), false))
	rt.router.DELETE("/photos/:photo_id", rt.auth(rt.wrap(rt.deletePhoto), false))
	rt.router.POST("/photos/:photo_id/comments", rt.auth(rt.wrap(rt.commentPhoto), false))
	rt.router.DELETE("/photos/:photo_id/comments/:comment_id", rt.auth(rt.wrap(rt.uncommentPhoto), false))
	rt.router.PUT("/photos/:photo_id/likes/:user_id", rt.auth(rt.wrap(rt.likePhoto), true))
	rt.router.DELETE("/photos/:photo_id/likes/:user_id", rt.auth(rt.wrap(rt.unlikePhoto), true))
	rt.router.GET("/users", rt.auth(rt.wrap(rt.searchUser), false))
	rt.router.GET("/users/:user_id", rt.auth(rt.wrap(rt.getUserProfile), false))
	rt.router.PATCH("/users/:user_id", rt.auth(rt.wrap(rt.setMyUserName), true))
	rt.router.GET("/users/:user_id/stream", rt.auth(rt.wrap(rt.getMyStream), false))
	rt.router.POST("/users/:user_id/photos", rt.auth(rt.wrap(rt.uploadPhoto), true))
	rt.router.GET("/users/:user_id/followers", rt.auth(rt.wrap(rt.getFollowers), false))
	rt.router.GET("/users/:user_id/following", rt.auth(rt.wrap(rt.getFollowing), false))
	rt.router.GET("/users/:user_id/banned", rt.auth(rt.wrap(rt.getBanned), true))
	rt.router.PUT("/users/:user_id/following/:target_user_id", rt.auth(rt.wrap(rt.followUser), true))
	rt.router.DELETE("/users/:user_id/following/:target_user_id", rt.auth(rt.wrap(rt.unfollowUser), true))
	rt.router.PUT("/users/:user_id/banned/:target_user_id", rt.auth(rt.wrap(rt.banUser), true))
	rt.router.DELETE("/users/:user_id/banned/:target_user_id", rt.auth(rt.wrap(rt.unbanUser), true))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
