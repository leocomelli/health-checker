package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/leocomelli/health-checker/core"
	db "github.com/leocomelli/health-checker/database"
	"github.com/leocomelli/health-checker/ping"
	"github.com/sirupsen/logrus"
)

func main() {

	services, err := core.LoadServices()
	if err != nil {
		logrus.Fatal(err)
	}

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("health", *services)
			return next(c)
		}
	})

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET},
	}))

	e.GET("/ping", ping.Check)
	e.GET("/database", db.Check)
	e.GET("/database/:sid", db.Check)

	e.Logger.Fatal(e.Start(":8080"))
}
