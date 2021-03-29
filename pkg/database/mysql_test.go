package database

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewDB(t *testing.T) {
	db, err := NewDB(MysqlCfg{
		Username: "mock",
		Password: "mocl",
		Host:     "mock",
		Port:     "3306",
		Name:     "mock",
		Driver:   "mysql",
	})
	assert.Nil(t, err)
	assert.NotNil(t, db)
	assert.NotNil(t, db.DB())
	assert.Nil(t, db.Close())
}

func TestNewDB_Error(t *testing.T) {
	db, err := NewDB(MysqlCfg{
		Username: "mock",
		Password: "mocl",
		Host:     "mock",
		Port:     "3306",
		Name:     "mock",
		Driver:   "error",
	})
	assert.NotNil(t, err)
	assert.Nil(t, db)
}
