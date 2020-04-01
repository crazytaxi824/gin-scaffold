package action

import (
	"net/http"
	"time"

	"src/global"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// var connPool map[string]interface{}
//
// func init() {
// 	connPool = make(map[string]interface{})
// }

const (
	buffSize = 1024
	timeout  = 5
)

var upgrader = websocket.Upgrader{
	HandshakeTimeout: timeout * time.Second,
	ReadBufferSize:   buffSize,
	WriteBufferSize:  buffSize,

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
