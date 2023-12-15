package clients

import (
	"context"
	"labs/service-mesh/dashboard/configs"
	custerror "labs/service-mesh/helper/error"
	custhttp "labs/service-mesh/helper/http"
	"labs/service-mesh/helper/logger"
	modelslocations "labs/service-mesh/helper/models/locations"
	"net/url"
	"time"

	fastshot "github.com/opus-domini/fast-shot"
	"go.uber.org/zap"
)

type LocationsServiceClient struct {
	httpClient *fastshot.Client
}

func NewLocationsServiceClient(ctx context.Context) *LocationsServiceClient {
	cfg := configs.Get()
	return &LocationsServiceClient{
		httpClient: custhttp.NewHttpClient(
			ctx,
			custhttp.
				WithBaseUrl(cfg.GetServices().Locations),
			custhttp.WithTimeout(time.Second*2),
		),
	}
}

func (c *LocationsServiceClient) GetCurrentLocation(ctx context.Context, req *modelslocations.GetCurrentLocationRequest) (*modelslocations.GetCurrentLocationResponse, error) {
	path, err := c.buildGetCurrentLocationPath(req)
	if err != nil {
		logger.SDebug("GetCurrentLocation: buildGetCurrentLocationPath",
			zap.Error(err))
		return nil, err
	}

	apiResp, err := c.httpClient.GET(path).
		Context().Set(ctx).
		Send()

	if err != nil {
		logger.SDebug("GetCurrentLocation: GetCurrentLocation",
			zap.Error(err))
		return nil, err
	}

	if apiResp.Is4xxClientError() {
		return nil, custerror.ErrorInvalidArgument
	}

	if apiResp.Is5xxServerError() {
		return nil, custerror.ErrorInternal
	}

	var resp modelslocations.GetCurrentLocationResponse
	if err := custhttp.JSONResponse(&apiResp, &resp); err != nil {
		logger.SDebug("GetCurrentLocation: JSONResponse",
			zap.Error(err))
		return nil, err
	}

	return &resp, nil
}

func (c *LocationsServiceClient) buildGetCurrentLocationPath(req *modelslocations.GetCurrentLocationRequest) (string, error) {
	if req.IP == "" {
		return "", custerror.ErrorInvalidArgument
	}

	p, _ := url.ParseRequestURI("/api/locate")
	q := p.Query()
	q.Add("ip", req.IP)
	p.RawQuery = q.Encode()

	return p.String(), nil
}

func (c *LocationsServiceClient) GetCurrentTime(ctx context.Context, req *modelslocations.GetCurrentTimeRequest) (*modelslocations.GetCurrentTimeResponse, error) {
	path, err := c.buildGetCurrentTime(req)
	if err != nil {
		logger.SDebug("GetCurrentTime: buildGetCurrentLocationPath",
			zap.Error(err))
		return nil, err
	}

	apiResp, err := c.httpClient.GET(path).
		Context().Set(ctx).
		Send()

	if err != nil {
		logger.SDebug("GetCurrentTime: GetCurrentLocation",
			zap.Error(err))
		return nil, err
	}

	if apiResp.Is4xxClientError() {
		return nil, custerror.ErrorInvalidArgument
	}

	if apiResp.Is5xxServerError() {
		return nil, custerror.ErrorInternal
	}

	var resp modelslocations.GetCurrentTimeResponse
	if err := custhttp.JSONResponse(&apiResp, &resp); err != nil {
		logger.SDebug("GetCurrentTime: JSONResponse",
			zap.Error(err))
		return nil, err
	}

	return &resp, nil
}

func (c *LocationsServiceClient) buildGetCurrentTime(req *modelslocations.GetCurrentTimeRequest) (string, error) {
	if req.IP == "" {
		return "", custerror.ErrorInvalidArgument
	}

	p, _ := url.ParseRequestURI("/api/time")
	q := p.Query()
	q.Add("ip", req.IP)
	p.RawQuery = q.Encode()

	return p.String(), nil
}
