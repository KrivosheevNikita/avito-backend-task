package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	RequestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "pvz_requests_total",
			Help: "Общее количество запросов",
		},
		[]string{"method", "path"},
	)

	ResponseDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "pvz_response_duration_seconds",
			Help:    "Длительность обработки запроса",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)
)

var (
	PVZCreated = promauto.NewCounter(prometheus.CounterOpts{
		Name: "pvz_created_total",
		Help: "Количество созданных ПВЗ",
	})

	ReceptionsCreated = promauto.NewCounter(prometheus.CounterOpts{
		Name: "receptions_created_total",
		Help: "Количество созданных приемок заказов",
	})

	ProductsAdded = promauto.NewCounter(prometheus.CounterOpts{
		Name: "items_added_total",
		Help: "Количество добавленных товаров",
	})
)

func Init() {
	prometheus.MustRegister(RequestCount)
	prometheus.MustRegister(ResponseDuration)
}
