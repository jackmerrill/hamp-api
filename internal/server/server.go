package server

import (
	"net/http"

	_ "github.com/jackmerrill/hamp-api/docs"
	menu "github.com/jackmerrill/hamp-api/internal/server/routes/dining"
	social "github.com/jackmerrill/hamp-api/internal/server/routes/social"
	util "github.com/jackmerrill/hamp-api/internal/server/routes/utilities"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Hamp API
// @version 1.0
// @description An API for various Hampshire College things.

// @contact.name Jack Merrill
// @contact.url https://jackmerrill.com
// @contact.email me@jackmerrill.com

// @license.name MPL 2.0
// @license.url https://www.mozilla.org/en-US/MPL/2.0/

// @host api.hamp.sh
// @BasePath /api
func Start() error {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.GET("/", HealthCheck)

	e.GET("/docs/*", echoSwagger.WrapHandler)

	api := e.Group("/api")

	api.GET("/", HealthCheck)

	socialGroup := api.Group("/social")
	utilitiesGroup := api.Group("/utilities")
	dining := api.Group("/dining")

	overHeard := socialGroup.Group("/overheard")
	overHeard.POST("/generate", social.GeneratePost)

	laundry := utilitiesGroup.Group("/laundry")
	util.InitLaundry()
	laundry.GET("/dakin", util.GetDakin)
	laundry.GET("/dakin/machines/:machine", util.GetDakinMachine)
	laundry.GET("/dakin/live", util.GetDakinLive)

	laundry.GET("/merrill", util.GetMerrill)
	laundry.GET("/merrill/machines/:machine", util.GetMerrillMachine)
	laundry.GET("/merrill/live", util.GetMerrillLive)

	laundry.GET("/enfield", util.GetEnfield)
	laundry.GET("/enfield/machines/:machine", util.GetEnfieldMachine)
	laundry.GET("/enfield/live", util.GetEnfieldLive)

	laundry.GET("/prescott", util.GetPrescott)
	laundry.GET("/prescott/machines/:machine", util.GetPrescottMachine)
	laundry.GET("/prescott/live", util.GetPrescottLive)

	dining.GET("/menu", menu.GetMenu)
	dining.GET("/menu/today", menu.GetTodaysMenu)

	return e.Start(":1323")
}

// HealthCheck godoc
// @Summary Show the status of server.
// @Description get the status of server.
// @Tags root
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router / [get]
func HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": "Server is up and running",
	})
}
