package monitoring

import (
	"github.com/jmoiron/sqlx"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/blacksponge/vertica-prometheus-exporter/db"
)

// PrometheusMetric maps a struct to a Prometheus valid map.
type PrometheusMetric interface {
	ToMetric() map[string]float64
}

type VerticaCollector struct {
	server *db.Server
	up prometheus.Gauge
}

func NewPrometheusMetrics(db sqlx.DB) []PrometheusMetric {
	var metrics []PrometheusMetric

	//for _, state := range NewNodeState(&db) {
	//	metrics = append(metrics, state)
	//}
	//for _, rejection := range NewPoolRejections(&db) {
	//	metrics = append(metrics, rejection)
	//}
	//for _, queryRequest := range NewQueryRequests(&db) {
	//	metrics = append(metrics, queryRequest)
	//}
	//for _, usage := range NewPoolUsage(&db) {
	//	metrics = append(metrics, usage)
	//}
	//metrics = append(metrics, NewVerticaSystem(&db))
	metrics = append(metrics, NewLicenseCompliance(&db))

	return metrics
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

// ToSnakeCase converts all string values to snake case.
func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func NewDesc(name string, labels []string) *prometheus.Desc {
	return prometheus.NewDesc(
		prometheus.BuildFQName("vertica", "", ToSnakeCase(name)),
		"", labels, nil,
	)
}

func NewVerticaCollect(server *db.Server) VerticaCollector {
	c := VerticaCollector {
		server: server,
	}
	c.up = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace:   "vertica",
		Name:        "up",
		Help:        "Whether the last scrape of metrics from Vertica was able to connect to the server (1 for yes, 0 for no).",
	})
	return c
}

func (collector VerticaCollector) Describe(ch chan<- *prometheus.Desc) {

}

func (collector VerticaCollector) Collect(ch chan<- prometheus.Metric) {
	db, err := collector.server.GetDB()

	if err != nil {
		log.Errorf("Could not connect to vertica: %v", err)
		collector.up.Set(0)
	} else {
		collector.up.Set(1)
	}

	ch <- collector.up

	if db != nil {

		for _, state := range NewNodeState(db) {
			state.Collect(ch)
		}
		for _, system := range NewVerticaSystem(db) {
			system.Collect(ch)
		}
		for _, queryRequest := range NewQueryRequests(db) {
			queryRequest.Collect(ch)
		}
		for _, usage := range NewPoolUsage(db) {
			usage.Collect(ch)
		}
		for _, rejection := range NewPoolRejections(db) {
			rejection.Collect(ch)
		}
		NewLicenseCompliance(db).Collect(ch)
	}
}
