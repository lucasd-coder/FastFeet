package monitor

import (
	"fmt"
	"regexp"

	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	OK    = "Ok"
	ERROR = "Error"
)

func RegisterSrvMetrics() *grpcprom.ServerMetrics {
	srvMetrics := grpcprom.NewServerMetrics(
		grpcprom.WithServerHandlingTimeHistogram(
			grpcprom.WithHistogramBuckets([]float64{0.001, 0.01, 0.1, 0.3, 0.6, 1, 3, 6, 9, 20, 30, 60, 90, 120}),
		),
	)
	return srvMetrics
}

type Metrics interface {
	IncHits(status, queueName string)
	ObserveResponseTime(status, queueName string, observeTime float64)
}

type PrometheusMetrics struct {
	HitsTotal prometheus.Counter
	Hits      *prometheus.CounterVec
	Times     *prometheus.HistogramVec
}

func CreateMetrics(name string, reg *prometheus.Registry) (Metrics, error) {
	var metr PrometheusMetrics
	metr.HitsTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: extractLabelName(name) + "_hits_total",
	})
	if err := reg.Register(metr.HitsTotal); err != nil {
		return nil, fmt.Errorf("prometheus.Register: %w", err)
	}
	metr.Hits = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: extractLabelName(name) + "_hits",
		},
		[]string{"status", "queueName"},
	)
	if err := reg.Register(metr.Hits); err != nil {
		return nil, fmt.Errorf("prometheus.Register: %w", err)
	}
	metr.Times = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: extractLabelName(name) + "_times",
		},
		[]string{"status", "queueName"},
	)
	if err := reg.Register(metr.Times); err != nil {
		return nil, fmt.Errorf("prometheus.Register: %w", err)
	}

	return &metr, nil
}

func (metr *PrometheusMetrics) IncHits(status, queueName string) {
	metr.HitsTotal.Inc()
	metr.Hits.WithLabelValues(status, extractLabelName(queueName)).Inc()
}

func (metr *PrometheusMetrics) ObserveResponseTime(status, queueName string, observeTime float64) {
	metr.Times.WithLabelValues(status, extractLabelName(queueName)).Observe(observeTime)
}

func extractLabelName(input string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9_]`)
	return re.ReplaceAllString(input, "_")
}
