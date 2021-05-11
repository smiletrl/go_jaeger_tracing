package constants

const (
	// tracing endpoint
	TracingEndpoint string = "http://jaeger-collector.istio-system.svc.cluster.local:14268/api/traces"

	// this one should be in env var in prod env
	Stage string = "dev"
)
