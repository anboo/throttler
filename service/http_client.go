package service

import (
	"context"
	"net/http"
	"time"

	"github.com/anboo/throttler/service/rate_limiter"
	"github.com/anboo/throttler/service/storage"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type HttpClient struct {
	httpClient  *http.Client
	rateLimiter rate_limiter.RateLimiter
	logger      *zerolog.Logger
}

func NewHttpClient(rateLimiter rate_limiter.RateLimiter, logger *zerolog.Logger) *HttpClient {
	return &HttpClient{
		httpClient:  http.DefaultClient,
		rateLimiter: rateLimiter,
		logger:      logger,
	}
}

func (c *HttpClient) Request(ctx context.Context, r storage.Request) error {
	start := time.Now()

	err := c.rateLimiter.Wait(ctx)
	if err != nil {
		return errors.Wrap(err, "http client wait rate limiter")
	}

	c.logger.Info().Str("request", r.ID).Str("rate_limit_wait", time.Since(start).String()).Msg("call request")

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
