Go gRPC Opentracing Instrumentation
===================================

[![GoDoc](https://godoc.org/github.com/devoxel/otgrpc?status.svg)](https://godoc.org/github.com/devoxel/otgrpc)

As opposed to using ClientInterceptors, use GRPC's StatsHandlers, as done in [ocgrpc]("https://github.com/census-instrumentation/opencensus-go/tree/master/plugin/ocgrpc").

See GitHub's fork information for history. This fork specifically cuts down on logging (to avoid stream spans growing too large).

Usage
-----

Client side:

```go
tracer := // Tracer implementation

th := otgrpc.NewTraceHandler(tracer)
conn, err := grpc.Dial(address, grpc.WithStatsHandler(th))
```

Server side:

```go
tracer := // Tracer implementation

th := otgrpc.NewTraceHandler(tracer)
server := grpc.NewServer(grpc.StatsHandler(th))
```

### Options

Limit tracing to methods of your choosing

```go
tf := func(method string) bool {
    return method == "/my.svc/my.method"
}

th := otgrpc.NewTraceHandler(tracer, orgrpc.WithTraceEnabledFunc(tf))
```
