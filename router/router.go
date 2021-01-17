package router

import (
	"encoding/json"

	routing "github.com/fasthttp/router"
	"github.com/fasthttp/websocket"
	"github.com/valyala/fasthttp"
	"github.com/zoommix/fasthttp_template/resources/users"
	"github.com/zoommix/fasthttp_template/utils"
)

const (
	// CLOSE connection action key
	CLOSE = "close"
	// HELLO ..
	HELLO = "hello"
)

// New ...
func New() *routing.Router {
	router := routing.New()

	router.NotFound = NotFoundHandler

	withTiming := Middlewares{RequestInfo, TimingMW}
	withAuth := append(withTiming, AuthMW)

	router.ANY("/ws", wsRouter)

	// ./static/somefile_in_static_folder
	router.GET("/static/{file}", fasthttp.FSHandler(utils.GetPWD(), 0))

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

var upgrader = websocket.FastHTTPUpgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var ps = &utils.PubSub{}

func wsRouter(ctx *fasthttp.RequestCtx) {
	upgrader.CheckOrigin = func(ctx *fasthttp.RequestCtx) bool {
		return true
	}

	user := fetchWsUser(ctx)

	if user == nil {
		utils.LogInfo("Unable to authorize user")
		return
	}

	err := upgrader.Upgrade(ctx, func(conn *websocket.Conn) {
		client := utils.NewClient(conn)
		client.UserID = user.ID

		// add this client into the list
		ps.AddClient(client)

		for {
			_, payload, err := conn.ReadMessage()
			utils.LogInfo(string(payload))

			if err != nil {
				utils.LogError(err)
				ps.RemoveClient(client)
				utils.LogInfo("total clients and subscriptions ", len(ps.Clients), len(ps.Subscriptions))
				return
			}

			m := utils.NewMessage()

			if err = json.Unmarshal(payload, &m); err != nil {
				client.SendJSONError(m.Action, err.Error())
				utils.LogError(err)
				return
			}

			switch m.Action {

			case CLOSE:
				utils.LogInfo("Disconnecting client: ", client.ID)
				return

			case HELLO:
				users.HelloHandler(&client)

			default:
				break
			}
		}
	})

	if err != nil {
		utils.LogError(err)
	}
}

func fetchWsUser(ctx *fasthttp.RequestCtx) *users.User {
	token := string(ctx.QueryArgs().Peek("token"))

	if len(ctx.Request.Header.Peek("Authorization")) == 0 {
		ctx.Request.Header.Add("Authorization", token)
	}

	return fetchUserByCtx(ctx)
}
