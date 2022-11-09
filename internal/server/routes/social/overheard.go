package routes

import (
	"fmt"
	"log"

	"github.com/fogleman/gg"
	"github.com/labstack/echo/v4"
)

type PostBody struct {
	Location string `json:"location" form:"location" query:"location"`
	Content  string `json:"content" form:"content" query:"content"`
}

func GeneratePost(c echo.Context) error {
	body := new(PostBody)

	if err := c.Bind(body); err != nil {
		return err
	}

	if body.Content == "" {
		return c.String(400, "Content is required")
	}

	if body.Location == "" {
		body.Location = "Location Unknown"
	}

	img := gg.NewContext(1080, 1080)

	// Draw a white background
	img.SetHexColor("#FFFFFF")
	img.DrawRectangle(0, 0, 1080, 1080)
	img.Fill()

	img.SetHexColor("#000000")

	if err := img.LoadFontFace("./fonts/Roboto-BoldItalic.ttf", 40); err != nil {
		log.Println(err)
		return err
	}
	img.DrawString(fmt.Sprintf("Overheard in %s", body.Location), 100, 500)

	if err := img.LoadFontFace("./fonts/Roboto-Regular.ttf", 30); err != nil {
		log.Println(err)
		return err
	}

	img.DrawStringWrapped(body.Content, 100, 560, 0, 0, 960, 1.5, gg.AlignLeft)

	img.EncodePNG(c.Response().Writer)

	return nil
}
