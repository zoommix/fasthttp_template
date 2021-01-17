package router

import (
	"encoding/json"
	"log"

	routing "github.com/fasthttp/router"
	"github.com/fasthttp/websocket"
	"github.com/valyala/fasthttp"
	"github.com/zoommix/fasthttp_template/resources/users"
	"github.com/zoommix/fasthttp_template/utils"
)

const (
	// CLOSE connection action key
	CLOSE = "close"
)

// New ...
func New() *routing.Router {
	router := routing.New()

	router.NotFound = NotFoundHandler

	withTiming := Middlewares{RequestInfo, TimingMW}
	withAuth := append(withTiming, AuthMW)

	router.ANY("/ws", wsRouter)

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

	err := upgrader.Upgrade(ctx, func(conn *websocket.Conn) {
		client := utils.NewClient()
		client.Connection = conn
		client.User = users.User{}

		// add this client into the list
		ps.AddClient(client)

		for {
			_, payload, err := conn.ReadMessage()
			utils.LogInfo(string(payload))

			if err != nil {
				log.Println("Something went wrong", err)

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
				ps.Publish(m.Topic, m.Message, &client)
				return

			default:
				break
			}
		}
	})

	if err != nil {
		utils.LogError(err)
	}
}
