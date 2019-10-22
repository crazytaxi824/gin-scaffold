package action

import (
	"net/http"
	"src/global"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var connPool map[string]interface{}

func init() {
	connPool = make(map[string]interface{})
}

var upgrader = websocket.Upgrader{
	HandshakeTimeout: 5 * time.Second,
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,

	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WebsocketTest(ctx *gin.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		global.Logger.Error("upgrade: " + err.Error())
		return
	}
	defer conn.Close()

	for {
		// 接收消息
		_, msg, err := conn.ReadMessage()
		if err != nil {
			global.Logger.Error(err.Error())
			return
		}

		err = conn.WriteMessage(1, msg)
		if err != nil {
			global.Logger.Error(err.Error())
			return
		}
	}
}
