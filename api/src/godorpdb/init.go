package godorpdb

import (
	"time"

	"github.com/jinzhu/gorm"
	// importing for postgres driver
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// DB is database pointer to postgres orm.
var DB *gorm.DB
var err error

// Post structure for posts table.
type Post struct {
	gorm.Model
	Author  string
	Message string
}

func addDatabase(dbname string) error {
	// create database with dbname, won't do anything if db already exists
	DB.Exec("CREATE DATABASE " + dbname)

	// connect to newly created DB (now has dbname param)
	connectionParams := "dbname=" + dbname + " user=docker password=docker sslmode=disable host=db"
	DB, err = gorm.Open("postgres", connectionParams)
	if err != nil {
		return err
	}

	return nil
}

// InitDb function to init postgres database.
func InitDb() (*gorm.DB, error) {
	// set up DB connection and then attempt to connect 5 times over 25 seconds
	connectionParams := "user=docker password=docker sslmode=disable host=db"
	for i := 0; i < 5; i++ {
		DB, err = gorm.Open("postgres", connectionParams) // gorm checks Ping on Open
		if err == nil {
			break
		}
		time.Sleep(5 * time.Second)
	}

	if err != nil {
		return DB, err
	}

	// create table if it does not exist
	if !DB.HasTable(&Post{}) {
		DB.CreateTable(&Post{})
	}

	testPost := Post{Author: "Dorper", Message: "GoDoRP is Dope"}
	DB.Create(&testPost)

	return DB, err
}
