package ping

import (
	"net/http"
	"time"

	"github.com/labstack/echo"
)

type Response struct {
	Text     string    `json:"text"`
	DateTime time.Time `json:"dateTime"`
}

func Check(c echo.Context) error {
	return c.JSON(http.StatusCreated, Response{"pong", time.Now()})
}
