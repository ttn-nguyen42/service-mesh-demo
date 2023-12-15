package publicapi

import (
	"github.com/gofiber/fiber/v2"
)

func ServiceRegistration() func(app *fiber.App) {
	return func(app *fiber.App) {
		apiGroup := app.Group("api")
		apiGroup.Get("timezones", GetListAreas)
		apiGroup.Get("locate", GetCurrentLocation)
		apiGroup.Get("time", GetCurrentTime)
	}
}
