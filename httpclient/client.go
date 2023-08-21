package httpclient

import (
	"context"
	"github.com/cross-space-official/common/consts"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"time"
)

type AuthPayload struct {
	ProfileID string
	ApiKey    string
}

type clientFactory struct {
	headers map[string]func(ctx context.Context) string
}

type XSpaceHttpClient struct {
	client  *resty.Client
	headers map[string]func(ctx context.Context) string
}

var factory = &clientFactory{
	headers: make(map[string]func(ctx context.Context) string),
}

func GetHttpClientFactory() *clientFactory {
	return factory
}

func (cf *clientFactory) WithHeader(key string, accessor func(ctx context.Context) string) *clientFactory {
	cf.headers[key] = accessor
	return factory
}

func (cf *clientFactory) Build(baseURL string, auth AuthPayload, headers map[string]string) *XSpaceHttpClient {
	client := resty.New()

	client.
		SetBaseURL(baseURL).
		SetHeaders(headers).
		SetHeader("Content-Type", "application/json").
		SetRetryCount(3).
		SetRetryWaitTime(100 * time.Millisecond).
		SetRetryMaxWaitTime(1 * time.Second).
		SetLogger(logrus.New())

	if len(auth.ApiKey) > 0 {
		client.SetHeader(consts.ApiKeyHeader, auth.ApiKey)
	}
	if len(auth.ProfileID) > 0 {
		client.SetHeader(consts.ProfileIDHeader, auth.ProfileID)
	}

	return &XSpaceHttpClient{client: client, headers: cf.headers}
}

func (client *XSpaceHttpClient) BuildRequest(c context.Context) *resty.Request {
	request := client.client.R()
	if c == nil {
		return request
	}

	for key, accessor := range client.headers {
		request.SetHeader(key, accessor(c))
	}

	return request
}
