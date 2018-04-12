// package otgrpc provides Opentracing instrumentation for gRPC services
package otgrpc

import (
	context "golang.org/x/net/context"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"google.golang.org/grpc/stats"
)

const (
	FailFastKey = "failfast"
	ClientKey   = "client"
)

var (
	GRPCComponentTag = opentracing.Tag{Key: string(ext.Component), Value: "gRPC"}
)

// TraceHandler is an implementation of grpc.StatsHandler that provides built in trace handling.
// Use NewTraceHandler to create one.
type TraceHandler struct {
	tracer   opentracing.Tracer
	opts     *options
	disabled bool
}

// NewTraceHandler creates a gRPC stats.Handler instance that instruments RPCs with Opentracing trace contexts
func NewTraceHandler(tracer opentracing.Tracer, o ...Option) *TraceHandler {
	th := &TraceHandler{
		tracer: tracer,
		opts:   newOptions(o...),
	}
	if _, ok := tracer.(opentracing.NoopTracer); ok {
		th.disabled = true
	}
	return th
}

// TagRPC is called when the RPC begins
func (th *TraceHandler) TagRPC(ctx context.Context, tagInfo *stats.RPCTagInfo) context.Context {
	if th.disabled || !th.opts.traceEnabledFunc(tagInfo.FullMethodName) {
		return ctx
	}

	spanCtx, err := extractSpanContext(th.tracer, ctx)
	if err != nil {
		return ctx
	}
	span := th.tracer.StartSpan(tagInfo.FullMethodName, opentracing.FollowsFrom(spanCtx), GRPCComponentTag)
	newCtx, _ := injectSpanToMetadata(th.tracer, span, ctx)
	return opentracing.ContextWithSpan(newCtx, span)
}

// HandleRPC is a catch all for all types of events that can happen during a stream.
func (th *TraceHandler) HandleRPC(ctx context.Context, s stats.RPCStats) {
	if th.disabled {
		return
	}
	span := opentracing.SpanFromContext(ctx)
	if span == nil {
		return
	}

	switch t := s.(type) {
	case *stats.Begin:
		span.LogFields(
			log.Bool(ClientKey, t.Client),
			log.Bool(FailFastKey, t.FailFast),
		)
	case *stats.End:
		if t.IsClient() {
			span.SetTag(string(ext.SpanKind), ext.SpanKindRPCClientEnum)
		} else {
			span.SetTag(string(ext.SpanKind), ext.SpanKindRPCServerEnum)
		}

		if t.Error != nil {
			span.SetTag(string(ext.Error), true)
			span.LogFields(log.Error(t.Error))
		}
		span.Finish()
	}
}

func (th *TraceHandler) TagConn(ctx context.Context, tagInfo *stats.ConnTagInfo) context.Context {
	return ctx
}

func (th *TraceHandler) HandleConn(ctx context.Context, s stats.ConnStats) {}
