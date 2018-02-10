package responder

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/snsinfu/web-skttt/domain"
)

var (
	errorMap = map[error]int{
		domain.ErrTopicNotFound: http.StatusNotFound,
		domain.ErrInvalidKey:    http.StatusUnauthorized,
	}
)

// Error maps given error to an HTTP response. nil is mapped to 200 OK.
func Error(c echo.Context, err error) error {
	if err == nil {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"error": nil,
		})
	}

	status, ok := errorMap[err]
	if !ok {
		status = http.StatusInternalServerError
	}

	return c.JSON(status, map[string]string{
		"error": err.Error(),
	})
}

// Data sends given data in JSON.
func Data(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"error": nil,
		"data":  data,
	})
}
