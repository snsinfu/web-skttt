package action

import (
	"io/ioutil"

	"github.com/labstack/echo"
	"github.com/snsinfu/web-skttt/domain"
	"github.com/snsinfu/web-skttt/non"
	"github.com/snsinfu/web-skttt/responder"
)

// POST /request
func (act *Action) PostRequest(c echo.Context) error {
	topic, key, err := act.dom.CreateTopic()
	if err != nil {
		return responder.Error(c, err)
	}

	return responder.Data(c, map[string]string{
		"topic": topic,
		"key":   key,
	})
}

// POST /:topic/:key
// POST /:topic
func (act *Action) PostTopic(c echo.Context) error {
	topic := c.Param("topic")
	key, _ := non.Empty(c.Param("key"), c.Request().Header.Get("X-API-Key"))

	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return responder.Error(c, err)
	}

	msg := domain.Message{
		Data: body,
	}
	return responder.Error(c, act.dom.PostToTopic(topic, key, msg))
}

// GET /:topic
func (act *Action) GetTopic(c echo.Context) error {
	sub, err := act.dom.Subscribe(c.Param("topic"))
	if err != nil {
		return responder.Error(c, err)
	}
	defer act.dom.Unsubscribe(sub)

	return responder.Stream(c, sub.Read())
}
