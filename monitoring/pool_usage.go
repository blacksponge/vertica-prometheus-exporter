package monitoring

import (
	"log"

	"github.com/jmoiron/sqlx"

	"github.com/prometheus/client_golang/prometheus"
)

// PoolUsage shows gneral resource pool usage stats.
type PoolUsage struct {
	NodeName                string `db:"node_name"`
	PoolName                string `db:"pool_name"`
	MemoryInUseKB           float64    `db:"memory_inuse_kb"`
	GeneralMemoryBorrowedKB float64    `db:"general_memory_borrowed_kb"`
	RunningQueryCount       float64    `db:"running_query_count"`
}

// NewPoolUsage returns a list of pool usage stats.
func NewPoolUsage(db *sqlx.DB) []PoolUsage {
	sql := `
	SELECT
		node_name,
		pool_name,
		memory_inuse_kb,
		general_memory_borrowed_kb,
		running_query_count
	FROM resource_pool_status;
	`

	usage := []PoolUsage{}
	err := db.Select(&usage, sql)
	if err != nil {
		log.Fatal(err)
	}

	return usage
}

func (usage PoolUsage) Collect(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(
		NewDesc("pool_memory_inuse_kb", []string { "node_name", "pool_name" }),
		prometheus.GaugeValue,
		usage.MemoryInUseKB,
		usage.NodeName,
		usage.PoolName,
	)

	ch <- prometheus.MustNewConstMetric(
		NewDesc("pool_memory_borrowed_kb", []string { "node_name", "pool_name" }),
		prometheus.GaugeValue,
		usage.GeneralMemoryBorrowedKB,
		usage.NodeName,
		usage.PoolName,
	)

	ch <- prometheus.MustNewConstMetric(
		NewDesc("pool_running_query_count", []string { "node_name", "pool_name" }),
		prometheus.GaugeValue,
		usage.RunningQueryCount,
		usage.NodeName,
		usage.PoolName,
	)
}
