package middleware

import (
	"context"
	"net/http"
	"pvz-service/internal/metrics"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
)

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		method := r.Method

		timer := prometheus.NewTimer(metrics.ResponseDuration.WithLabelValues(method, path))
		defer timer.ObserveDuration()

		metrics.RequestCount.WithLabelValues(method, path).Inc()

		next.ServeHTTP(w, r)
	})
}

func GrpcMetrics(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	start := time.Now()

	method := info.FullMethod
	path := method

	resp, err := handler(ctx, req)

	metrics.RequestCount.WithLabelValues(method, path).Inc()
	metrics.ResponseDuration.WithLabelValues(method, path).Observe(time.Since(start).Seconds())

	return resp, err
}
