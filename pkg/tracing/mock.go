package tracing

import (
	"context"
	"io"

	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"
)

type mockProvider struct{}

func NewMockProvider() Provider {
	return mockProvider{}
}

func (m mockProvider) SetupTracer(serviceName, endpoint, stage string) (io.Closer, error) {
	return nil, nil
}

func (m mockProvider) Middleware() echo.MiddlewareFunc {
	return nil
}

func (m mockProvider) StartSpan(c context.Context, operationName string) (opentracing.Span, context.Context) {
	return nil, context.Background()
}

func (p mockProvider) FinishSpan(span opentracing.Span) {
}
