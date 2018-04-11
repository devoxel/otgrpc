package otgrpc

import (
	"context"

	"google.golang.org/grpc/metadata"

	"github.com/opentracing/opentracing-go"
)

func extractSpanContext(tracer opentracing.Tracer, ctx context.Context) (opentracing.SpanContext, error) {
	sc := spanContextFromContext(ctx)
	if sc != nil {
		return sc, nil
	}
	return extractSpanContextFromMetadata(tracer, ctx)
}

func spanContextFromContext(ctx context.Context) opentracing.SpanContext {
	if parentSpan := opentracing.SpanFromContext(ctx); parentSpan != nil {
		return parentSpan.Context()
	}
	return nil
}

func injectSpanToMetadata(tracer opentracing.Tracer, span opentracing.Span, ctx context.Context) (context.Context, error) {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md = metadata.New(nil)
	}
	if err := tracer.Inject(span.Context(), opentracing.HTTPHeaders, NewMetadataReaderWriter(md)); err != nil {
		return ctx, err
	}
	return metadata.NewOutgoingContext(ctx, md), nil
}

func extractSpanContextFromMetadata(tracer opentracing.Tracer, ctx context.Context) (opentracing.SpanContext, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.New(nil)
	}
	return tracer.Extract(opentracing.HTTPHeaders, NewMetadataReaderWriter(md))
}
