package monitoring

import (
	"log"

	"github.com/jmoiron/sqlx"

	"github.com/prometheus/client_golang/prometheus"
)

// QueryRequest lists query performance metrics on the username level.
type QueryRequest struct {
	UserName          string `db:"user_name"`
	RequestDurationMS int    `db:"request_duration_ms"`
	MemoryAcquiredMB  int    `db:"memory_acquired_mb"`
}

// NewQueryRequests returns query performance for all users.
func NewQueryRequests(db *sqlx.DB) []QueryRequest {
	sql := `
	SELECT
		user_name,
		SUM(COALESCE(request_duration_ms,0))::INT request_duration_ms,
		SUM(COALESCE(memory_acquired_mb,0))::INT memory_acquired_mb
	FROM v_monitor.query_requests
	GROUP BY user_name;
	`

	queryRequests := []QueryRequest{}
	err := db.Select(&queryRequests, sql)
	if err != nil {
		log.Fatal(err)
	}

	return queryRequests
}

func (qr QueryRequest) Collect(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(
		NewDesc("request_duration_ms", []string { "user_name" }),
		prometheus.GaugeValue,
		float64(qr.RequestDurationMS),
		qr.UserName,
	)

	ch <- prometheus.MustNewConstMetric(
		NewDesc("memory_acquired_mb", []string { "user_name" }),
		prometheus.GaugeValue,
		float64(qr.MemoryAcquiredMB),
		qr.UserName,
	)
}
