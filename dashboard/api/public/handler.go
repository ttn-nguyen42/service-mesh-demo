package publicapi

import (
	"labs/service-mesh/dashboard/biz/handlers"
	custerror "labs/service-mesh/helper/error"
	"labs/service-mesh/helper/logger"
	modelsdash "labs/service-mesh/helper/models/dashboard"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func GetDashboard(ctx *fiber.Ctx) error {
	ips := ctx.IPs()
	ip := ""
	if len(ips) > 0 {
		ip = ips[0]
	}
	logger.SInfo("GetDashboard: GetIP",
		zap.String("ip", ip))
	if len(ip) == 0 {
		return custerror.ErrorUnavailable
	}

	resp, err := handlers.GetHandlers().GetDashboardData(ctx.Context(),
		&modelsdash.GetDashboardDataRequest{
			IP: ip,
		})
	if err != nil {
		return err
	}

	return ctx.JSON(resp)
}
