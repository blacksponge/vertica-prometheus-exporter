package monitoring

import (
	"log"

	"github.com/fatih/structs"
	"github.com/jmoiron/sqlx"

	"github.com/prometheus/client_golang/prometheus"
)

// VerticaSystem shows important system values such as the epoch and fault tolerance levels.
type VerticaSystem struct {
	CurrentEpoch           float64 `db:"current_epoch"`
	AhmEpoch               float64 `db:"ahm_epoch"`
	LastGoodEpoch          float64 `db:"last_good_epoch"`
	RefreshEpoch           float64 `db:"refresh_epoch"`
	DesignedFaultTolerance float64 `db:"designed_fault_tolerance"`
	NodeCount              float64 `db:"node_count"`
	NodeDownCount          float64 `db:"node_down_count"`
	CurrentFaultTolerance  float64 `db:"current_fault_tolerance"`
	CatalogRevisionNumber  float64 `db:"catalog_revision_number"`
	WosUsedBytes           float64 `db:"wos_used_bytes"`
	WosRowCount            float64 `db:"wos_row_count"`
	RosUsedBytes           float64 `db:"ros_used_bytes"`
	RosRowCount            float64 `db:"ros_row_count"`
	TotalUsedBytes         float64 `db:"total_used_bytes"`
	TotalRowCount          float64 `db:"total_row_count"`
}

// NewVerticaSystem returns a new instance of VerticaSystem
func NewVerticaSystem(db *sqlx.DB) []VerticaSystem {
	sql := `SELECT * FROM system`

	system := []VerticaSystem{}
	err := db.Select(&system, sql)
	if err != nil {
		log.Fatal(err)
	}

	return system
}

func (sys VerticaSystem) Collect(ch chan<- prometheus.Metric) {
	for k, v := range structs.Map(sys) {
		ch <- prometheus.MustNewConstMetric(NewDesc(k, nil), prometheus.GaugeValue, v.(float64))
	}
}
