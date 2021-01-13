package router

import (
	"fmt"
	"strconv"
	"time"

	"github.com/valyala/fasthttp"
	"github.com/zoommix/fasthttp_template/resources/users"
	"github.com/zoommix/fasthttp_template/utils"
)

// Middleware ...
type Middleware func(ctx fasthttp.RequestHandler) fasthttp.RequestHandler

// Middlewares ..
type Middlewares []Middleware

// ApplyToHandler ...
func (ms Middlewares) ApplyToHandler(handler fasthttp.RequestHandler) fasthttp.RequestHandler {
	for i := range ms {
		idx := len(ms) - 1 - i
		handler = ms[idx](handler)
	}

	return handler
}

// RequestInfo ...
func RequestInfo(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		requestInfo := fmt.Sprintf("Started: [%s] %s?%s", string(ctx.Method()), string(ctx.Path()), ctx.QueryArgs())
		utils.LogNotice(requestInfo)
		next(ctx)
	}
}

// TimingMW ..
func TimingMW(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		start := time.Now()
		next(ctx)
		utils.LogNotice("Elapsed time: ", time.Since(start))
	}
}

// AuthMW ..
func AuthMW(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		u := fetchUser(ctx)

		if u != nil {
			ctx.SetUserValue("user", u)
			next(ctx)
			return
		}

		utils.RenderUnauthorized(ctx, "token invalid or missing")
	}
}

func fetchUser(ctx *fasthttp.RequestCtx) *users.User {
	token := string(ctx.Request.Header.Peek("Authorization"))
	u, err := users.DecodeJWT(token)

	if err != nil {
		utils.LogError("Unable to verify JWT token: ", err.Error())
	} else {
		utils.LogInfo("Logged user ID = " + strconv.Itoa(u.ID))
	}

	return u
}
