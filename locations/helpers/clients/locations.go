package dtclient

import (
	"context"
	custerror "labs/service-mesh/helper/error"
	custhttp "labs/service-mesh/helper/http"
	"labs/service-mesh/helper/logger"
	"net/url"

	fastshot "github.com/opus-domini/fast-shot"
	"go.uber.org/zap"
)

type LocationApiClient struct {
	httpClient *fastshot.Client
}

type GetCurrentLocationRequest struct {
	IP string `json:"ip"`
}

// 17031891
// https://ip-api.com/docs/api:json
type IPInfo struct {
	Query       string  `json:"query"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	City        string  `json:"city"`
	Latitude    float64 `json:"lat"`
	Longitude   float64 `json:"lon"`
	ISP         string  `json:"isp"`
	Mobile      bool    `json:"mobile"`
	Proxy       bool    `json:"proxy"`
	Hosting     bool    `json:"hosting"`
}

type GetCurrentLocationResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	IPInfo
}

func (c *LocationApiClient) GetCurrentLocation(ctx context.Context, req *GetCurrentLocationRequest) (*GetCurrentLocationResponse, error) {
	uri, err := c.buildGetCurrentLocationPath(req)
	if err != nil {
		logger.SInfo("GetCurrentLocation: buildGetCurrentLocationPath",
			zap.Error(err))
		return nil, err
	}

	apiResp, err := c.httpClient.GET(uri).
		Context().Set(ctx).
		Send()
	if err != nil {
		logger.SInfo("GetCurrentLocation: Send",
			zap.Error(err))
		return nil, err
	}

	if apiResp.Is4xxClientError() {
		return nil, custerror.ErrorInvalidArgument
	}

	if apiResp.Is5xxServerError() {
		return nil, custerror.ErrorInternal
	}

	var resp GetCurrentLocationResponse
	if err := custhttp.JSONResponse(&apiResp, &resp); err != nil {
		logger.SInfo("GetCurrentLocation: JSONResponse",
			zap.Error(err))
		return nil, err
	}

	return &resp, nil
}

func (c *LocationApiClient) buildGetCurrentLocationPath(req *GetCurrentLocationRequest) (string, error) {
	if req.IP == "" {
		return "", custerror.FormatInvalidArgument("missing IP")
	}

	path, err := url.ParseRequestURI("/json")
	if err != nil {
		return "", err
	}
	path = path.JoinPath(req.IP)
	q := path.Query()
	q.Add("fields", "17031891")
	path.RawQuery = q.Encode()

	return path.String(), nil
}
