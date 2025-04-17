package server

import (
	"net/http"

	"github.com/TerraDharitri/drt-go-chain-core/core/check"
	"github.com/TerraDharitri/drt-go-chain-core/marshal"
	"github.com/TerraDharitri/drt-go-chain/api/logs"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// NewGinHandler will create a gin handler
func NewGinHandler(marshaller marshal.Marshalizer) (*gin.Engine, error) {
	if check.IfNilReflect(marshaller) {
		return nil, errNilMarshaller
	}

	router := gin.Default()
	registerLoggerWsRoute(router, marshaller)

	return router, nil
}

func registerLoggerWsRoute(ws *gin.Engine, marshaller marshal.Marshalizer) {
	upgrader := websocket.Upgrader{}

	ws.GET("/log", func(c *gin.Context) {
		upgrader.CheckOrigin = func(r *http.Request) bool {
			return true
		}

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Error(err.Error())
			return
		}

		ls, err := logs.NewLogSender(marshaller, conn, log)
		if err != nil {
			log.Error(err.Error())
			return
		}

		ls.StartSendingBlocking()
	})
}
