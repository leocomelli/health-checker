package db

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/leocomelli/health-checker/core"
	_ "gopkg.in/goracle.v2"
)

const stmt = `SELECT INSTANCE_NAME, HOST_NAME, VERSION, STARTUP_TIME, STATUS FROM V$INSTANCE`

var errTimeout = errors.New("max execution time exceeded")

type Response struct {
	InstanceName string
	HostName     string
	Version      string
	StartupTime  time.Time
	Status       string
	DateTime     time.Time
	Error        error
}

func Check(c echo.Context) error {
	health := c.Get("health").(core.Health)
	dbServices := health.GetByType("database")

	responses := make(chan Response)

	for _, s := range dbServices {
		go func(s core.Service) {
			db, err := sql.Open("goracle", s.URL)
			if err != nil {
				responses <- Response{Error: err}
			} else {
				var i, h, v, ss string
				var t time.Time
				err = db.QueryRow(stmt).Scan(&i, &h, &v, &t, &ss)
				if err != nil {
					responses <- Response{Error: err}
				} else {
					responses <- Response{i, h, v, t, ss, time.Now(), nil}
				}
			}
		}(s)
	}

	var rs []Response
	for {
		select {
		case r := <-responses:
			rs = append(rs, r)
			if len(rs) == len(dbServices) {
				return c.JSON(http.StatusCreated, rs)
			}
		case <-time.After(60 * time.Second):
			rs = append(rs, Response{Error: errTimeout})
			return c.JSON(http.StatusCreated, rs)
		}
	}
}
