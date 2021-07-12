package main

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/smiletrl/go_jaeger_tracing/pkg/constants"
	"github.com/smiletrl/go_jaeger_tracing/pkg/tracing"
	rpcserver "github.com/smiletrl/go_jaeger_tracing/service.product/internal/rpc/server"
)

func main() {
	// echo instance
	e := echo.New()

	tracingProvider := tracing.NewProvider()
	closer, err := tracingProvider.SetupTracer("cart", constants.TracingEndpoint, constants.Stage)
	if err != nil {
		panic(err)
	}
	defer closer.Close()

	// middleware
	e.Use(tracingProvider.Middleware())

	group := e.Group("/api/v1")
	group.GET("", Get)

	// start grpc server
	go func() {
		err = rpcserver.Register()
		if err != nil {
			panic(err)
		}
	}()

	// start rest server
	if err := e.Start(":1324"); err != nil {
		panic(err)
	}
}

func Get(c echo.Context) error {
	return c.String(http.StatusOK, "ok")
}
