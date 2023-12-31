package publicapi

import (
	custerror "labs/service-mesh/helper/error"
	modelscities "labs/service-mesh/helper/models/cities"
	"labs/service-mesh/locations/biz/service"

	"github.com/gofiber/fiber/v2"
)

func GetCountries(ctx *fiber.Ctx) error {
	resp, err := service.
		GetCitiesService().
		GetCountries(ctx.Context())
	if err != nil {
		return err
	}

	return ctx.JSON(resp)
}

func GetCities(ctx *fiber.Ctx) error {
	c := ctx.Query("country")
	if len(c) == 0 {
		return custerror.ErrorInvalidArgument
	}
	resp, err := service.
		GetCitiesService().
		GetCities(ctx.Context(), &modelscities.GetCitiesRequest{
			Iso2: c,
		})
	if err != nil {
		return err
	}

	return ctx.JSON(resp)
}
