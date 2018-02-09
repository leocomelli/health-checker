package database

import (
	"database/sql"
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo"
	"github.com/leocomelli/health-checker/core"
	_ "gopkg.in/goracle.v2"
)

const (
	stmt = `SELECT INSTANCE_NAME, HOST_NAME, VERSION, STARTUP_TIME, STATUS FROM V$INSTANCE`
)

type Response struct {
	InstanceName string
	HostName     string
	Version      string
	StartupTime  time.Time
	Status       string
	Error        error
}

func Check(c echo.Context) error {
	health := c.Get("health").(core.Health)
	dbServices := health.GetByType("database")

	responses := make(chan Response)

	var wg sync.WaitGroup
	wg.Add(len(dbServices))

	for _, s := range dbServices {
		go func(s core.Service) {
			defer wg.Done()

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
					responses <- Response{i, h, v, t, ss, nil}
				}
			}
		}(s)
	}

	var rs []Response
	go func() {
		for r := range responses {
			rs = append(rs, r)
		}
	}()

	wg.Wait()

	return c.JSON(http.StatusCreated, rs)
}
