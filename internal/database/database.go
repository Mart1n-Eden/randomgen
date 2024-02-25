package database 

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/driver/postgres"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "postgres"
)

var (
	db *gorm.DB
)

func Init(entity interface{}) {
	dbinfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	var err error

	db,err = gorm.Open(postgres.Open(dbinfo),  &gorm.Config{})
	if  err != nil {
		panic(err)
	}

	err = db.AutoMigrate(entity)
	if err != nil {
		panic(err)
	}
}

func AddItem(entity interface{}) error {
	if res := db.Create(entity); res.Error != nil {
		return res.Error
	}

	return nil
}