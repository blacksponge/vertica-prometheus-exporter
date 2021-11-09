package monitoring

import (
	"log"

	"github.com/jmoiron/sqlx"

	"github.com/prometheus/client_golang/prometheus"
)

// NodeState contains information about each node in a Vertica cluster.
type NodeState struct {
	NodeID    string `db:"node_id"`
	NodeName  string `db:"node_name"`
	NodeState float64    `db:"node_state"`
}

// NewNodeState returns the status for each node in the Vertica cluster.
func NewNodeState(db *sqlx.DB) []NodeState {
	sql := `
	SELECT
		node_id,
		node_name,
		(node_state='UP')::INT node_state
	FROM v_catalog.nodes`

	nodeState := []NodeState{}
	err := db.Select(&nodeState, sql)
	if err != nil {
		log.Fatal(err)
	}

	return nodeState
}

func (ns NodeState) Collect(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(
		NewDesc("node_state", []string { "node_id", "node_name" }),
		prometheus.GaugeValue,
		ns.NodeState,
		ns.NodeID,
		ns.NodeName,
	)
}
