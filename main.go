package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	api "github.com/anboo/throttler/http"
	"github.com/anboo/throttler/resource"
	"github.com/anboo/throttler/service"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	l := zerolog.New(zerolog.NewConsoleWriter())

	res := resource.NewResources()
	res.Initialize()

	httpClient := service.NewHttpClient(res.RateLimiter, &l)

	queue := service.NewQueue(
		ctx,
		res.Env.IntervalCheckingNewRequests,
		httpClient,
		runtime.GOMAXPROCS(0),
		res.Storage,
		&l,
	)

	createRequestHandler := api.NewCreateRequestHandler(res.Storage)
	getRequestHandler := api.NewGetRequestHandler(res.Storage)

	r := gin.Default()
	r.POST("/request", createRequestHandler.Handler)
	r.GET("/request/:id", getRequestHandler.Handler)

	srv := &http.Server{Addr: ":8080", Handler: r}
	go func() {
		queue.Start(ctx)
	}()
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-ctx.Done()
	cancel()
	l.Warn().Msg("shutting down")

	timeoutCtx, c := context.WithTimeout(context.Background(), 10*time.Second)
	defer c()
	if err := srv.Shutdown(timeoutCtx); err != nil {
		l.Fatal().Err(err).Msg("server shutdown")
	}
	l.Warn().Msg("server exiting")
}
