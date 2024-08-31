package menu

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

// GetTomorrowsMenu godoc
// @Summary Get tomorrow's Dining menu
// @Description Get tomorrow's Dining menu
// @Tags dining
// @Accept json
// @Produce json
// @Success 200 {object} Meal
// @Router /dining/menu/tomorrow [get]
func GetTomorrowsMenu(c echo.Context) error {
	// Parse the website for the URL
	url, err := ParseWebsite()
	if err != nil {
		panic(err)
	}

	// Parse the URL for the CSV
	parsedURL, err := ParseURL(*url)
	if err != nil {
		panic(err)
	}

	// Parse the CSV for the Menu
	menu, err := ParseCSV(*parsedURL)
	if err != nil {
		panic(err)
	}

	if menu == nil {
		panic("menu is nil")
	}

	var todaysMenu Meal
	todaysDate := time.Now().Format("1/2/2006")

	for i, v := range *menu {
		if v.Date == todaysDate {
			todaysMenu = (*menu)[i+1]
			break
		}
	}

	return c.JSON(http.StatusOK, todaysMenu)
}
