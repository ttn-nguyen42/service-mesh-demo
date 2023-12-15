package publicapi

import (
	custerror "labs/service-mesh/helper/error"
	modelsweather "labs/service-mesh/helper/models/weather"
	"labs/service-mesh/weather/biz/handlers"

	"github.com/gofiber/fiber/v2"
)

func GetLocation(ctx *fiber.Ctx) error {
	query := ctx.Queries()
	city, found := query["city"]
	if !found {
		return custerror.ErrorInvalidArgument
	}
	resp, err := handlers.GetHandlers().
		GetLocation(ctx.Context(), &modelsweather.GetLocationRequest{
			City: city,
		})
	if err != nil {
		return err
	}
	return ctx.JSON(resp)
}

func GetCurrentWeather(ctx *fiber.Ctx) error {
	lat := ctx.QueryFloat("latitude")
	if lat == 0 {
		return custerror.ErrorInvalidArgument
	}
	lon := ctx.QueryFloat("longitude")
	if lon == 0 {
		return custerror.ErrorInvalidArgument
	}
	resp, err := handlers.
		GetHandlers().
		GetCurrentWeather(ctx.Context(), &modelsweather.GetCurrentWeatherRequest{
			Latitude:  float32(lat),
			Longitude: float32(lon),
		})
	if err != nil {
		return err
	}

	return ctx.JSON(resp)
}
