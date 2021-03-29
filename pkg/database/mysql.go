package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type mysqlConnect struct {
	db *sql.DB
}

type MysqlCfg struct {
	Username string
	Password string
	Host     string
	Port     string
	Name     string
	Driver   string
}

func (d *mysqlConnect) DB() *sql.DB {
	return d.db
}

func (d *mysqlConnect) Close() error {
	return d.db.Close()
}

func NewDB(cfg MysqlCfg) (Connector, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=UTC",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
	)
	db, err := sql.Open(cfg.Driver, dsn)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	return &mysqlConnect{
		db: db,
	}, nil
}
