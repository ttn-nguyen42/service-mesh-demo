package publicapi

import (
	custerror "labs/service-mesh/helper/error"
	modelslocations "labs/service-mesh/helper/models/locations"
	"labs/service-mesh/locations/biz/handlers"

	"github.com/gofiber/fiber/v2"
)

func GetListAreas(ctx *fiber.Ctx) error {
	resp, err := handlers.
		Location().
		GetListAreas(ctx.Context())
	if err != nil {
		return err
	}
	return ctx.JSON(resp)
}

func GetCurrentLocation(ctx *fiber.Ctx) error {
	ip := ctx.Query("ip")
	if len(ip) == 0 {
		return custerror.ErrorInvalidArgument
	}
	resp, err := handlers.
		Location().
		GetCurrentLocation(ctx.Context(),
			&modelslocations.GetCurrentLocationRequest{
				IP: ip,
			})
	if err != nil {
		return err
	}
	return ctx.JSON(resp)
}

func GetCurrentTime(ctx *fiber.Ctx) error {
	ip := ctx.Query("ip")
	if len(ip) == 0 {
		return custerror.ErrorInvalidArgument
	}
	resp, err := handlers.
		Location().
		GetCurrentTime(ctx.Context(),
			&modelslocations.GetCurrentTimeRequest{
				IP: ip,
			})
	if err != nil {
		return err
	}
	return ctx.JSON(resp)
}
