package db

import (
	_ "github.com/vertica/vertica-sql-go"
	"github.com/jmoiron/sqlx"
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
