package monitoring

import (
	"log"
	"time"
	"regexp"
	"strconv"

	"github.com/jmoiron/sqlx"

	"github.com/prometheus/client_golang/prometheus"
)

type LicenseCompliance struct {
	Utilization			float64
	AuditTime			time.Time
	NodeCount			uint64
	NodeLimit			uint64
	RawDataSize			float64
	RawDataConfidence	float64
	LicenseSize			float64
}

type complianceStatus struct {
	compliance_status	string
}

var licenseSizeRegex = regexp.MustCompile(`License Size : ([0-9\.]+)TB`)
var utilizationRegex = regexp.MustCompile(`Utilization  : ([0-9]+)%`)
var nodeCountRegex   = regexp.MustCompile(`Node count : ([0-9]+)`)
var nodeLimitRegex   = regexp.MustCompile(`License Node limit : ([0-9]+)`)
var rawDataRegex     = regexp.MustCompile(`Raw Data Size: ([0-9\.]+)TB \+/- ([0-9\.]+)TB`)
var auditTimeRegex   = regexp.MustCompile(`Audit Time   : (.+)`)

// NewNodeState returns the status for each node in the Vertica cluster.
func NewLicenseCompliance(db *sqlx.DB) LicenseCompliance {
	sql := `SELECT get_compliance_status();`

	complianceStatusOut := []string{}
	licenseCompliance :=  LicenseCompliance{}
	err := db.Select(&complianceStatusOut, sql)
	if err != nil {
		log.Fatal(err)
	}

	licenseCompliance.LicenseSize, err = strconv.ParseFloat(licenseSizeRegex.FindStringSubmatch(complianceStatusOut[0])[1], 64)
	if err != nil {
		log.Fatal(err)
	}

	licenseCompliance.Utilization, err = strconv.ParseFloat(utilizationRegex.FindStringSubmatch(complianceStatusOut[0])[1], 64)
	if err != nil {
		log.Fatal(err)
	}
	licenseCompliance.Utilization = licenseCompliance.Utilization / 100

	licenseCompliance.NodeCount, err = strconv.ParseUint(nodeCountRegex.FindStringSubmatch(complianceStatusOut[0])[1], 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	licenseCompliance.NodeLimit, err = strconv.ParseUint(nodeLimitRegex.FindStringSubmatch(complianceStatusOut[0])[1], 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	rawData := rawDataRegex.FindStringSubmatch(complianceStatusOut[0])
	licenseCompliance.RawDataSize, err = strconv.ParseFloat(rawData[1], 64)
	if err != nil {
		log.Fatal(err)
	}

	licenseCompliance.RawDataConfidence, err = strconv.ParseFloat(rawData[2], 64)
	if err != nil {
		log.Fatal(err)
	}

	licenseCompliance.AuditTime, err = time.Parse("2006-01-02 15:04:05.999999999-07", auditTimeRegex.FindStringSubmatch(complianceStatusOut[0])[1])
	if err != nil {
		log.Fatal(err)
	}

	return licenseCompliance
}

// ToMetric converts NodeState to a Map.
func (lc LicenseCompliance) ToMetric() map[string]float64 {
	metrics := map[string]float64{}

	metrics["vertica_license_compliance_license_size"] = lc.LicenseSize
	metrics["vertica_license_compliance_utilization"] = lc.Utilization
	metrics["vertica_license_compliance_node_count"] = float64(lc.NodeCount)
	metrics["vertica_license_compliance_node_limmit"] = float64(lc.NodeLimit)
	metrics["vertica_license_compliance_raw_data_size"] = lc.RawDataSize
	metrics["vertica_license_compliance_raw_data_confidence"] = lc.RawDataConfidence
	metrics["vertica_license_compliance_raw_audit_tiime"] = float64(lc.AuditTime.Unix())

	return metrics
}

func (lc LicenseCompliance) Collect(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(
		NewDesc("license_compliance_license_size", nil),
		prometheus.GaugeValue,
		lc.LicenseSize,
	)

	ch <- prometheus.MustNewConstMetric(
		NewDesc("license_compliance_utilization", nil),
		prometheus.GaugeValue,
		lc.Utilization,
	)

	ch <- prometheus.MustNewConstMetric(
		NewDesc("license_compliance_node_count", nil),
		prometheus.GaugeValue,
		float64(lc.NodeCount),
	)

	ch <- prometheus.MustNewConstMetric(
		NewDesc("license_compliance_node_limmit", nil),
		prometheus.GaugeValue,
		float64(lc.NodeLimit),
	)

	ch <- prometheus.MustNewConstMetric(
		NewDesc("license_compliance_raw_data_size", nil),
		prometheus.GaugeValue,
		lc.RawDataSize,
	)

	ch <- prometheus.MustNewConstMetric(
		NewDesc("license_compliance_raw_data_confidence", nil),
		prometheus.GaugeValue,
		lc.RawDataConfidence,
	)

	ch <- prometheus.MustNewConstMetric(
		NewDesc("license_compliance_raw_audit_tiime", nil),
		prometheus.GaugeValue,
		float64(lc.AuditTime.Unix()),
	)
}
