package dtclient

import (
	"context"
	custerror "labs/service-mesh/helper/error"
	custhttp "labs/service-mesh/helper/http"
	"labs/service-mesh/helper/logger"
	modelscities "labs/service-mesh/helper/models/cities"
	"path"

	fastshot "github.com/opus-domini/fast-shot"
	"go.uber.org/zap"
)

type CitiesApiClient struct {
	httpClient *fastshot.Client
}

func (c *CitiesApiClient) GetCountries(ctx context.Context) (modelscities.GetListCountriesResponse, error) {
	p := c.buildGetCountriesPath()

	resp, err := c.httpClient.GET(p).
		Context().Set(ctx).
		Send()
	if err != nil {
		logger.SInfo("GetCountries.Send", zap.Error(err))
		return nil, err
	}

	if resp.Is4xxClientError() {
		return nil, custerror.ErrorInvalidArgument
	}

	if resp.Is5xxServerError() {
		return nil, custerror.ErrorInternal
	}

	results := modelscities.GetListCountriesResponse{}
	if err := custhttp.JSONResponse(&resp, &results); err != nil {
		logger.SInfo("GetCountries.JSONResponse",
			zap.Error(err))
		return nil, err
	}

	return results, nil
}

func (c *CitiesApiClient) GetCities(ctx context.Context, req *GetCitiesRequest) (modelscities.GetListCitiesResponse, error) {
	p, err := c.buildGetCitiesPath(req)
	if err != nil {
		logger.SInfo("GetCities.Send", zap.Error(err))
		return nil, err
	}

	resp, err := c.httpClient.GET(p).
		Context().Set(ctx).
		Send()
	if err != nil {
		logger.SInfo("GetCountries.Send", zap.Error(err))
		return nil, err
	}

	if resp.Is4xxClientError() {
		return nil, custerror.ErrorInvalidArgument
	}

	if resp.Is5xxServerError() {
		return nil, custerror.ErrorInternal
	}

	results := modelscities.GetListCitiesResponse{}
	if err := custhttp.JSONResponse(&resp, &results); err != nil {
		logger.SInfo("GetCountries.JSONResponse",
			zap.Error(err))
		return nil, err
	}

	return results, nil
}

func (c *CitiesApiClient) buildGetCountriesPath() string {
	return "/v1/countries"
}

func (c *CitiesApiClient) buildGetCitiesPath(req *GetCitiesRequest) (string, error) {
	if req.Iso2 == "" {
		return "", custerror.ErrorInvalidArgument
	}

	return path.Join("/v1/countries", req.Iso2, "states"), nil
}
