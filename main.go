package main

import (
	"fmt"
	"net/http"

	"github.com/namsral/flag"
	log "github.com/sirupsen/logrus"

    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/blacksponge/vertica-prometheus-exporter/monitoring"
	"github.com/blacksponge/vertica-prometheus-exporter/db"
)


func main() {
	location := flag.String("location", "/metrics", "Metrics path")
	listen := flag.String("listen", "0.0.0.0:8080", "Address to listen on")
	dbUsername := flag.String("db_user", "dbadmin", "Vertica username")
	dbPassword := flag.String("db_password", "dbadmin", "Vertica password")
	dbHost := flag.String("db_host", "localhost", "Vertica hostname")
	dbPort := flag.Int("db_port", 5433, "Vertica port")
	dbName := flag.String("db_name", "vertica", "Vertica database name")
	flag.Parse()

	connString := fmt.Sprintf("vertica://%v:%v@%v:%v/%v", *dbUsername, *dbPassword, *dbHost, *dbPort, *dbName)
	server := db.NewServer(connString)
	collector := monitoring.NewVerticaCollect(server)
	prometheus.MustRegister(collector)

	log.Infof("starting serving metrics at %s%s", *listen, *location)
	http.Handle(*location, promhttp.Handler())
	err := http.ListenAndServe(*listen, nil)
	log.Fatal(err)
}
