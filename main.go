package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/anboo/throttler/service"
	"github.com/caarlos0/env/v6"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"golang.org/x/time/rate"
)

type Config struct {
	DatabaseDSN              string        `env:"DB_DSN"`
	RequestsLimitPerInterval time.Duration `env:"REQUESTS_LIMIT_PER_INTERVAL"`
	RequestsLimit            int           `env:"REQUESTS_INTERVAL"`
}

var config Config

func main() {
	ctx := context.Background()
	l := zerolog.New(os.Stdout)

	l.Info().Msg("parse env config start")
	if err := env.Parse(&config); err != nil {
		l.Fatal().Err(err).Msg("try parse env")
	}
	l.Info().Msg("parse env config done")

	rateLimiter := rate.NewLimiter(
		rate.Every(config.RequestsLimitPerInterval),
		config.RequestsLimit,
	)
	httpClient := service.NewHttpClient(rateLimiter)

	httpClient.Request(ctx)

	r := gin.Default()
	r.POST("/task", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.Run()
}
