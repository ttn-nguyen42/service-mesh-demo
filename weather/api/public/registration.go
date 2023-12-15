package publicapi

import (
	"github.com/gofiber/fiber/v2"
)

func ServiceRegistration() func(app *fiber.App) {
	return func(app *fiber.App) {
		apiGroup := app.Group("api")
		apiGroup.Get("location", GetLocation)
		apiGroup.Get("weather", GetCurrentWeather)
	}
}
