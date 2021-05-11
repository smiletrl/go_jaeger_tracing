package main

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/smiletrl/go_jaeger_tracing/pkg/constants"
	"github.com/smiletrl/go_jaeger_tracing/pkg/tracing"
	productClient "github.com/smiletrl/go_jaeger_tracing/service.product/external"
	"net/http"
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

	// product grpc client.
	pclient, err := productClient.NewClient(constants.GrpcProductHost, tracingProvider)
	if err != nil {
		panic(err)
	}
	productProxy := product{pclient}

	r := resource{tracingProvider, productProxy}
	group := e.Group("/api/v1")
	group.POST("/cart", r.Create)

	// start server
	if err = e.Start(":1325"); err != nil {
		panic(err)
	}
}

/* ----- resource ------ */

type resource struct {
	tracing      tracing.Provider
	productProxy ProductProxy
}

type createRequest struct {
	Quantity int    `json:"quantity"`
	SkuID    string `json:"sku_id"`
}

func (r resource) Create(c echo.Context) error {
	ctx := c.Request().Context()
	req := new(createRequest)

	if err := c.Bind(req); err != nil {
		panic(err)
	}

	// get product sku stock
	stock, err := r.productProxy.GetSkuStock(ctx, req.SkuID)
	if err != nil {
		panic(err)
	}

	if stock < req.Quantity {
		panic(err)
	}

	// save the cart data into database.
	if err = SaveCart(ctx, req.SkuID, req.Quantity, r.tracing); err != nil {
		panic(err)
	}

	return c.String(http.StatusOK, "ok")
}

/* --  product proxy --- */
type ProductProxy interface {
	GetSkuStock(c context.Context, skuID string) (int, error)
}

// product proxy
type product struct {
	client productClient.Client
}

func (p product) GetSkuStock(c context.Context, skuID string) (int, error) {
	return p.client.GetSkuStock(c, skuID)
}

/* ---- database --- */
func SaveCart(c context.Context, skuID string, quantity int, tracing tracing.Provider) error {
	// save it to redis/postgres
	span, ctx := tracing.StartSpan(c, "cart item create")
	defer span.Finish()

	// pass ctx to db related functions.
	fmt.Println(ctx)
	return nil
}
