package dbcon

import (
	"fmt"
	"log"
	"os"
	"vaccinationDrive/conf"

	"github.com/go-pg/pg"
)

var db *pg.DB

//Connect database
func Connect() {
	dbCon := pg.Connect(&pg.Options{
		Addr:            conf.Cfg.DB_ADDRESS,
		User:            conf.Cfg.DB_USERNAME,
		Password:        conf.Cfg.DB_PASSWORD,
		Database:        conf.Cfg.DB_NAME,
		PoolSize:        conf.Cfg.Max_Connection_Pool_Size,
		ApplicationName: conf.Cfg.APP_NAME,
	})

	db = dbCon
	log.Printf("Connected successfully")

	_, err := db.Exec("SELECT 1")
	if err != nil {
		fmt.Println("PostgreSQL is down")
		log.Fatalf("Unable to connect Postgres Database: %v\n", err)
		os.Exit(1)
	}

	db.AddQueryHook(dbLogger{})
}

//Get db connection
func Get() *pg.DB {
	return db
}

//Close db connection
func Close() {
	err := db.Close()

	if err != nil {
		log.Printf("Closing DB err", err)
	}
	log.Printf("DB closed")
}

type dbLogger struct{}

func (d dbLogger) BeforeQuery(q *pg.QueryEvent) {
}

func (d dbLogger) AfterQuery(q *pg.QueryEvent) {
	// fmt.Println(q.FormattedQuery())

}
