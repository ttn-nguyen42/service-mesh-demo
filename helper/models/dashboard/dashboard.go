package modelsdash

import (
	modelslocations "labs/service-mesh/helper/models/locations"
	modelsweather "labs/service-mesh/helper/models/weather"
)

type GetDashboardDataRequest struct {
	IP string `json:"ip"`
}

type GetDashboardDataResponse struct {
	IP          string                                      `json:"ip"`
	CurrentTime *modelslocations.GetCurrentTimeResponse     `json:"current_time"`
	Location    *modelslocations.GetCurrentLocationResponse `json:"location"`
	Weather     *modelsweather.GetCurrentWeatherResponse    `json:"weather"`
}
