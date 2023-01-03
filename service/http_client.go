package service

import (
	"context"
	"net/http"

	"github.com/pkg/errors"
	"golang.org/x/time/rate"
)

type HttpClient struct {
	httpClient  *http.Client
	rateLimiter *rate.Limiter
}

func NewHttpClient(rateLimiter *rate.Limiter) *HttpClient {
	return &HttpClient{
		httpClient:  http.DefaultClient,
		rateLimiter: rateLimiter,
	}
}

func (c *HttpClient) Request(ctx context.Context) error {
	err := c.rateLimiter.Wait(ctx)
	if err != nil {
		return errors.Wrap(err, "")
	}

	req, err := http.NewRequest("GET", "google.com", nil)
	if err != nil {
		return errors.Wrap(err, "cannot create request")
	}

	_, err = c.httpClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "request error")
	}

	return nil
}
