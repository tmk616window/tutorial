package util

import (
	"database/sql"
	"log"

	"github.com/DATA-DOG/go-txdb"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	gormDB *gorm.DB
)

type Suite struct {
	suite.Suite
}

func (s *Suite) SetupDB() {
	for _, d := range sql.Drivers() {
		if d == "txdb" {
			return
		}
	}
	txdb.Register(
		"txdb",
		"postgres",
		"postgres://user:password@postgresql:5432/db?sslmode=disable",
	)
}

func (s *Suite) CloseDB() {
	if gormDB != nil {
		sqlDB, _ := gormDB.DB()
		sqlDB.Close()
	}
}

func (s *Suite) DB() *gorm.DB {
	gormDB, err := gorm.Open(postgres.Open("postgres://user:password@postgresql:5432/db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	return gormDB.Debug()
}
