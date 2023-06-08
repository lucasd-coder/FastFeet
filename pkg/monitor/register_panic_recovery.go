package monitor

import (
	"context"
	"runtime/debug"

	"github.com/lucasd-coder/router-service/pkg/logger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func RegisterGrpcPanicRecoveryHandler(ctx context.Context, reg prometheus.Registerer) func(p any) (err error) {
	log := logger.FromContext(ctx)

	panicsTotal := promauto.With(reg).NewCounter(prometheus.CounterOpts{
		Name: "grpc_req_panics_recovered_total",
		Help: "Total number of gRPC requests recovered from internal panic.",
	})

	return func(p any) (err error) {
		panicsTotal.Inc()
		log.Error("recovered from panic", "panic", p, "stack", debug.Stack())
		return status.Errorf(codes.Internal, "%s", p)
	}
}
