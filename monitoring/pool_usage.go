package monitoring

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

// PoolUsage shows gneral resource pool usage stats.
type PoolUsage struct {
	NodeName                string `db:"node_name"`
	PoolName                string `db:"pool_name"`
	MemoryInUseKB           int    `db:"memory_inuse_kb"`
	GeneralMemoryBorrowedKB int    `db:"general_memory_borrowed_kb"`
	RunningQueryCount       int    `db:"running_query_count"`
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

// ToMetric converts PoolUsage to a Map.
func (usage PoolUsage) ToMetric() map[string]float64 {
	metrics := map[string]float64{}

	node := fmt.Sprintf("node_name=%q", usage.NodeName)
	pool := fmt.Sprintf("pool_name=%q", usage.PoolName)
	metrics[fmt.Sprintf("vertica_pool_memory_inuse_kb{%s, %s}", node, pool)] = float64(usage.MemoryInUseKB)
	metrics[fmt.Sprintf("vertica_pool_memory_borrowed_kb{%s, %s}", node, pool)] = float64(usage.GeneralMemoryBorrowedKB)
	metrics[fmt.Sprintf("vertica_pool_running_query_count{%s, %s}", node, pool)] = float64(usage.RunningQueryCount)

	return metrics
}
