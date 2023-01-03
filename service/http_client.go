package service

import (
	"context"
	"net/http"

	"github.com/anboo/throttler/service/storage"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"golang.org/x/time/rate"
)

type HttpClient struct {
	httpClient  *http.Client
	rateLimiter *rate.Limiter
	logger      *zerolog.Logger
}

func NewHttpClient(rateLimiter *rate.Limiter, logger *zerolog.Logger) *HttpClient {
	return &HttpClient{
		httpClient:  http.DefaultClient,
		rateLimiter: rateLimiter,
		logger:      logger,
	}
}

func (c *HttpClient) Request(ctx context.Context, r storage.Request) error {
	if !c.rateLimiter.Allow() {
		c.logger.Warn().Str("request", r.ID).Msg("rate limit reached maximum number of requests")
	} else {
		c.logger.Info().Str("request", r.ID).Msg("call request")
	}

	err := c.rateLimiter.Wait(ctx)
	if err != nil {
		return errors.Wrap(err, "http client wait rate limiter")
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
