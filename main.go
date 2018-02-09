package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"
	"github.com/leocomelli/health-checker/core"
	"github.com/leocomelli/health-checker/database"
	"github.com/leocomelli/health-checker/ping"
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
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))

	e.GET("/ping", ping.Check)
	e.GET("/database", database.Check)

	e.Run(standard.New(":8080"))
}
