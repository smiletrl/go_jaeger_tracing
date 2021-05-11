package external

import (
	"context"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/smiletrl/go_jaeger_tracing/pkg/constants"
	"github.com/smiletrl/go_jaeger_tracing/pkg/tracing"
	pb "github.com/smiletrl/go_jaeger_tracing/service.product/internal/rpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"time"
)

type Client interface {
	// Get sku stock
	GetSkuStock(ctx context.Context, skuID string) (stock int, err error)
}

type client struct {
	grpc    pb.ProductClient
	tracing tracing.Provider
}

func NewClient(endpoint string, tracingProvider tracing.Provider) (Client, error) {
	conn, err := newConnectionClient(endpoint)
	if err != nil {
		return nil, err
	}
	return client{
		grpc:    conn,
		tracing: tracingProvider,
	}, nil
}

func newConnectionClient(endpoint string) (client pb.ProductClient, err error) {
	var address = endpoint + constants.GrpcPort

	var kacp = keepalive.ClientParameters{
		Time:                10 * time.Second, // send pings every 10 seconds if there is no activity
		Timeout:             time.Second,      // wait 1 second for ping ack before considering the connection dead
		PermitWithoutStream: true,             // send pings even without active streams
	}

	// only allow maximum 1 second connection.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, address,
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithKeepaliveParams(kacp),
		grpc.WithStreamInterceptor(grpc_middleware.ChainStreamClient(
			grpc_opentracing.StreamClientInterceptor(),
		)),
		grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(
			grpc_opentracing.UnaryClientInterceptor(),
		)),
	)
	if err != nil {
		return nil, err
	}
	return pb.NewProductClient(conn), nil
}

func (c client) GetSkuStock(ctx context.Context, skuID string) (stock int, err error) {
	pbstock, err := c.grpc.GetSkuStock(ctx, &pb.SkuID{Value: skuID})
	if err != nil {
		return stock, err
	}

	return int(pbstock.Value), nil
}
