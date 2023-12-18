package dtclient

import (
	"context"
	"fmt"
	custerror "labs/service-mesh/helper/error"
	custhttp "labs/service-mesh/helper/http"
	"labs/service-mesh/helper/logger"
	"strings"

	fastshot "github.com/opus-domini/fast-shot"
	"go.uber.org/zap"
)

type DatetimeApiClient struct {
	httpClient *fastshot.Client
}

type Timezone struct {
	Location string `json:"location"`
	Area     string `json:"area"`
	Region   string `json:"region,omitempty"`
}

type Timezones []Timezone

type CurrentTimeRequest struct {
	Area     string `json:"area,omitempty"`
	Location string `json:"location,omitempty"`
	Region   string `json:"region,omitempty"`
	IP       string `json:"ip,omitempty"`
}

type GetTimezonesRequest struct {
	Area string `json:"area,omitempty"`
}

type CurrentTime struct {
	Abbreviation string      `json:"abbreviation,omitempty"`
	ClientIP     string      `json:"client_ip,omitempty"`
	DateTime     string      `json:"datetime,omitempty"`
	DayOfWeek    int         `json:"day_of_week,omitempty"`
	DayOfYear    int         `json:"day_of_year,omitempty"`
	DST          bool        `json:"dst,omitempty"`
	DSTFrom      interface{} `json:"dst_from,omitempty"`
	DSTOffset    int         `json:"dst_offset,omitempty"`
	DSTUntil     interface{} `json:"dst_until,omitempty"`
	RawOffset    int         `json:"raw_offset,omitempty"`
	Timezone     string      `json:"timezone,omitempty"`
	Unixtime     int         `json:"unixtime,omitempty"`
	UTCDateTime  string      `json:"utc_datetime,omitempty"`
	UTCOffset    string      `json:"utc_offset,omitempty"`
	WeekNumber   int         `json:"week_number,omitempty"`
}

func (c *DatetimeApiClient) GetCurrentTime(ctx context.Context, req CurrentTimeRequest) (*CurrentTime, error) {
	path, err := c.buildCurrentTimePath(&req)
	if err != nil {
		logger.SInfo("GetCurrentTime.buildCurrentTimePath",
			zap.Error(err))
		return nil, err
	}
	reqBuilder := c.httpClient.
		GET(path).
		Context().Set(ctx)

	resp, err := reqBuilder.Send()
	if err != nil {
		logger.SInfo("GetCurrentTime.Send",
			zap.Error(err))
		return nil, err
	}

	if resp.Is4xxClientError() {
		return nil, custerror.ErrorInvalidArgument
	}

	if resp.Is5xxServerError() {
		return nil, custerror.ErrorInternal
	}

	var result CurrentTime
	if err := custhttp.JSONResponse(&resp, &result); err != nil {
		logger.SInfo("GetCurrentTime.JSONResponse",
			zap.Error(err))
		return nil, err
	}

	return &result, nil
}

func (c *DatetimeApiClient) GetTimezones(ctx context.Context, req GetTimezonesRequest) (Timezones, error) {
	path := c.buildTimezonePath(&req)

	resp, err := c.httpClient.
		GET(path).
		Context().Set(ctx).
		Send()

	if err != nil {
		logger.SInfo("GetTimezones.Send",
			zap.Error(err))
		return nil, err
	}

	if resp.Is4xxClientError() {
		return nil, custerror.ErrorInvalidArgument
	}

	if resp.Is5xxServerError() {
		return nil, custerror.ErrorInternal
	}

	rawResult := []string{}
	if err := custhttp.JSONResponse(&resp, &rawResult); err != nil {
		logger.SInfo("GetTimezones.JSONResponse",
			zap.Error(err))
		return nil, err
	}

	results := Timezones{}
	for _, z := range rawResult {
		results = append(results, c.parseTimezones(z))
	}

	return results, nil
}

func (c *DatetimeApiClient) buildCurrentTimePath(req *CurrentTimeRequest) (string, error) {
	if req.IP != "" {
		ipBasePath := fmt.Sprintf("/api/ip/%s", req.IP)
		return ipBasePath, nil
	}

	if req.Area == "" {
		return "", custerror.FormatInvalidArgument("empty area")
	}
	if req.Location == "" {
		return "", custerror.FormatInvalidArgument("empty location")
	}

	basePath := fmt.Sprintf("/api/timezone/%s/%s", req.Area, req.Location)
	if req.Region != "" {
		basePath = fmt.Sprintf("%s/%s", basePath, req.Region)
	}
	return basePath, nil
}

func (c *DatetimeApiClient) buildTimezonePath(req *GetTimezonesRequest) string {
	basePath := "/api/timezone"
	if req.Area != "" {
		basePath = fmt.Sprintf("%s/%s", basePath, req.Area)
	}
	return basePath
}

func (c *DatetimeApiClient) parseTimezones(zone string) Timezone {
	parts := strings.Split(zone, "/")
	tz := Timezone{}
	if len(parts) >= 2 {
		tz.Area = parts[0]
		tz.Location = parts[1]
	}
	if len(parts) == 3 {
		tz.Region = parts[2]
	}
	return tz
}
