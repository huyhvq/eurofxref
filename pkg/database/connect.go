package database

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type Connector interface {
	Close() error
	DB() *sql.DB
}
