# Go tracing
This is an example for jaeger tracing in go for micro services. It assumes [Jaeger](https://www.jaegertracing.io/) has been installed at kubernetes already, and you are already familiar with [Opentracing](https://opentracing.io/).

You might want to tweak the variables defined at `pkg/constants`, depending on your installment.

Key points:
- Two services: cart, product. Service product has set up grpc server and rest server. Service cart only sets up rest server. Service cart will send request to service product grpc server.
- Start the root span at `(p *provider) Middleware()`. This middleware will create a root span for every request.
- Pass the propagating span relation using `opentracing.StartSpanFromContext()`. So we will see a request span to service product grpc server has a parent request span to service cart, using the `context.Context`.
- Create custom span like `span, ctx := tracing.StartSpan(c, "cart item create")`. See `service.product/main.go - SaveCart()`.

You might want to check the full example at repository [micro ecommerce](https://github.com/smiletrl/micro_ecommerce).

See the demo result like

Jaeger UI
![Jaeger UI](https://raw.githubusercontent.com/smiletrl/go_jaeger_tracing/master/assets/Jeager%20UI.png)

## Commands

Test Post command

```
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"sku_id":"12","quantity":1}' \
  http://127.0.0.1:1325/cart
```

Proto Update

```
cd service.product/internal

protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    proto/product.proto
```