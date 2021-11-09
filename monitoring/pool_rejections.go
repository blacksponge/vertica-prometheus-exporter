package monitoring

import (
	"log"

	"github.com/jmoiron/sqlx"

	"github.com/prometheus/client_golang/prometheus"
)

// PoolRejection shows the amount of resource pool rejections per node.
type PoolRejection struct {
	Reason         string `db:"reason"`
	ResourceType   string `db:"resource_type"`
	NodeName       string `db:"node_name"`
	PoolName       string `db:"pool_name"`
	RejectionCount float64 `db:"rejection_count"`
}

// NewPoolRejections returns a list of resource pool rejections from Vertica.
func NewPoolRejections(db *sqlx.DB) []PoolRejection {
	sql := `
	SELECT
		node_name,
		pool_name,
		rejection_count,
		resource_type,
		reason
	FROM v_monitor.resource_rejections`

	rejections := []PoolRejection{}
	err := db.Select(&rejections, sql)
	if err != nil {
		log.Fatal(err)
	}

	return rejections
}

func (pr PoolRejection) Collect(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(
		NewDesc("pool_rejection_count", []string { "node_name", "pool_name", "resource_type", "reason" }),
		prometheus.GaugeValue,
		pr.RejectionCount,
		pr.NodeName,
		pr.PoolName,
		pr.ResourceType,
		pr.Reason,
	)
}
