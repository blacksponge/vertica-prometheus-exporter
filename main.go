package main

import (
	"fmt"
	"net/http"

	"github.com/namsral/flag"
	log "github.com/sirupsen/logrus"
	_ "github.com/vertica/vertica-sql-go"

	"github.com/jmoiron/sqlx"
	"github.com/blacksponge/vertica-prometheus-exporter/monitoring"
)


type Server struct {
	dataSourceName string
	db *sqlx.DB
}

func NewServer(dsn string) *Server {
	return &Server {
		dataSourceName: dsn,
	}
}

func (s *Server) GetDB() (*sqlx.DB, error) {
	if s.db == nil || s.db.Ping() != nil {
		db, err := sqlx.Connect("vertica", s.dataSourceName)
		if err != nil {
			return nil, err
		}
		s.db = db
	}
	return s.db, nil
}

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
	server := NewServer(connString)

	serveMetrics(*location, *listen, *server)
}

// Serve Vertica metrics at chosen address and url.
func serveMetrics(location, listen string, server Server) {

	h := func(w http.ResponseWriter, r *http.Request) {
		db, err := server.GetDB()
		if err != nil {
			log.Errorf("Could not connect to vertica: %v", err)
			fmt.Fprintf(w, "vertica_up 0\n")
			return
		}
		fmt.Fprintf(w, "vertica_up 1\n")
		metrics := monitoring.NewPrometheusMetrics(*db)
		for _, obj := range metrics {
			metric := obj.ToMetric()
			for key, value := range metric {
				fmt.Fprintf(w, "%s %f\n", key, value)
			}
		}
	}

	http.HandleFunc(location, h)
	log.Printf("starting serving metrics at %s%s", listen, location)
	err := http.ListenAndServe(listen, nil)
	log.Fatal(err)
}
