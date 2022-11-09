package server

import (
	routes "github.com/jackmerrill/hamp-api/internal/server/routes/social"
	"github.com/labstack/echo/v4"
)

func Start() error {
	e := echo.New()

	api := e.Group("/api")

	social := api.Group("/social")

	overHeard := social.Group("/overheard")
	overHeard.POST("/generate", routes.GeneratePost)

	return e.Start(":1323")
}
