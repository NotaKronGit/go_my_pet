package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	dbType = "postgres"
)

type PgConfig struct {
	Host     string `required:"true"`
	Port     uint64 `default:"5432"`
	User     string `required:"true"`
	Password string `required:"true"`
	DBName   string `required:"true"`
	SslMode  string `default:"disable"`
}

func InitDbConnection(cfg *PgConfig) (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=%v",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SslMode)
	db, err := sql.Open(dbType, psqlInfo)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
