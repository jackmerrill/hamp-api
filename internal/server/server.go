package server

import (
	social "github.com/jackmerrill/hamp-api/internal/server/routes/social"
	util "github.com/jackmerrill/hamp-api/internal/server/routes/utilities"
	"github.com/labstack/echo/v4"
)

func Start() error {
	e := echo.New()

	api := e.Group("/api")

	socialGroup := api.Group("/social")
	utilitiesGroup := api.Group("/utilities")

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

	return e.Start(":1323")
}
