package dbscripts

import (
	"log"
	"vaccinationDrive/dbcon"
	"vaccinationDrive/models"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

//InitDB initialize DB
func InitDB() {
	db := dbcon.Get()
	CreateTables(db)
}

//getModels function use to get all the masters from models
func getModels() []interface{} {
	return []interface{}{
		&models.User{},
	}
}

//CreateTables function is use to create master tables
func CreateTables(db *pg.DB) {
	for _, mod := range getModels() {
		if err := db.CreateTable(mod, &orm.CreateTableOptions{
			IfNotExists:   true,
			FKConstraints: true,
		}); err != nil {
			log.Printf("Error in creating tables, err:%s", err.Error())
		}
	}
}
