package responder

import (
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"github.com/snsinfu/web-skttt/domain"
)

var (
	upgrader = websocket.Upgrader{}
)

const (
	websocketTimeout = 5 * time.Second
)

func Stream(c echo.Context, msgCh <-chan domain.Message) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return Error(c, err)
	}
	defer ws.Close()

	for msg := range msgCh {
		if err := ws.WriteMessage(websocket.TextMessage, msg.Data); err != nil {
			c.Logger().Error(err)
			return nil
		}
	}

	err = ws.WriteControl(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
		time.Now().Add(websocketTimeout),
	)
	if err != nil {
		c.Logger().Error(err)
	}

	return nil
}
