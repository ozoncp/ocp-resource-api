package metrics

import "github.com/prometheus/client_golang/prometheus"

const (
	namespace = "ocp"
)

var reqCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "resource_cud_requests_total",
	}, []string{"request"},
)

func RegisterMetrics() {
	prometheus.DefaultRegisterer.MustRegister(reqCounter)
}

func IncReqCounter(req string) {
	reqCounter.With(prometheus.Labels{"request": req}).Inc()
}
