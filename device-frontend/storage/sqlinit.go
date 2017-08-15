package storage

import (
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

func DatabaseInit() (err error) {
	host := os.Getenv("MYSQL_SERVER_HOST")
	port := os.Getenv("MYSQL_SERVER_PORT")
	username := os.Getenv("MYSQL_LOGIN_USER")
	password := os.Getenv("MYSQL_USER_PASSWORD")

	var dataSourceName string = username + ":" + password + "@tcp(" + host + ":" + port + ")/device?parseTime=true"
	db, err = sqlx.Connect("mysql", dataSourceName)
	if err != nil {
		LogErr("Err from storage.sqlinit.DatabaseInit(): ", err)
		return err
	}
	db.SetMaxIdleConns(1000)
	db.SetMaxOpenConns(1000)

	if err = imeisInit(); err != nil {
		return err
	}
	return nil
}
