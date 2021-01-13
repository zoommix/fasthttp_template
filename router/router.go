package router

import (
	routing "github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"github.com/zoommix/fasthttp_template/resources/users"
	"github.com/zoommix/fasthttp_template/utils"
)

// New ...
func New() *routing.Router {
	router := routing.New()

	router.NotFound = NotFoundHandler

	withTiming := Middlewares{RequestInfo, TimingMW}
	withAuth := append(withTiming, AuthMW)

	usersGroup := router.Group("/users")
	usersGroup.GET("/{id}", withAuth.ApplyToHandler(users.ShowUsersHandler))
	usersGroup.POST("", withTiming.ApplyToHandler(users.CreateUserHandler))

	router.POST("/tokens", users.Authenticate)

	return router
}

// NotFoundHandler ...
func NotFoundHandler(ctx *fasthttp.RequestCtx) {
	utils.RenderNotFoundError(ctx, "route not found")
	return
}
