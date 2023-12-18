package custhttp

import (
	"context"
	"labs/service-mesh/helper/logger"
	"net/http"
	"time"

	"github.com/motemen/go-loghttp"
	"github.com/opus-domini/fast-shot"
	"go.uber.org/zap"
)

type Options struct {
	baseUrl string
	timeout time.Duration
	header  map[string]string
}

type ClientOptioner func(o *Options)

func WithBaseUrl(url string) ClientOptioner {
	return func(o *Options) {
		o.baseUrl = url
	}
}

func WithTimeout(dur time.Duration) ClientOptioner {
	return func(o *Options) {
		o.timeout = dur
	}
}

func WithHeader(key string, value string) ClientOptioner {
	return func(o *Options) {
		o.header[key] = value
	}
}

func NewHttpClient(ctx context.Context, opts ...ClientOptioner) *fastshot.Client {
	options := &Options{}
	for _, o := range opts {
		o(options)
	}

	fs := fastshot.NewClient(options.baseUrl).
		Config().SetFollowRedirects(true).
		Config().SetTimeout(options.timeout).
		Config().SetCustomTransport(getDefaultRoundTripper())

	fs.Header().AddAll(options.header)

	return fs.Build()

}

func getDefaultRoundTripper() http.RoundTripper {
	httpLogTransport := &loghttp.Transport{
		Transport: http.DefaultTransport,
		LogRequest: func(req *http.Request) {
			logger.SDebug("HTTP Request",
				zap.String("scheme", req.URL.Scheme),
				zap.String("hostname", req.URL.Hostname()),
				zap.String("port", req.URL.Port()),
				zap.String("path", req.URL.EscapedPath()),
				zap.String("queries", req.URL.RawQuery),
				zap.Any("header", req.Header),
			)
		},
		LogResponse: func(resp *http.Response) {
			status := resp.StatusCode
			logger.SDebug("HTTP Response",
				zap.Int("status", status),
				zap.Any("headers", resp.Header))
		},
	}

	return &DurationHttpTransport{
		parentTransport: httpLogTransport,
	}
}

type DurationHttpTransport struct {
	parentTransport *loghttp.Transport
}

func (t *DurationHttpTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	reqStart := time.Now()
	resp, err := t.parentTransport.RoundTrip(r)
	reqEnd := time.Now()

	logger.SInfo("DurationHttpTransport",
		zap.Duration("duration", reqEnd.Sub(reqStart)))

	return resp, err
}
