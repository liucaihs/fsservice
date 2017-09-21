package main

import (
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var mysqdb *sqlx.DB

func mysqInit() (err error) {

	os.Setenv("DB_CONNECTION", "mysql")
	os.Setenv("DB_HOST", "192.168.1.211")
	os.Setenv("DB_PORT", "3306")
	os.Setenv("DB_DATABASE", "csv_import")
	os.Setenv("DB_USERNAME", "root")
	os.Setenv("DB_PASSWORD", "secretpw")
	os.Setenv("ALERT_EMAIL", "1258877243@qq.com")

	dbdrive := os.Getenv("DB_CONNECTION")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_DATABASE")
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	var dataSourceName string = username + ":" + password + "@tcp(" + host + ":" + port + ")/" + dbname + "?parseTime=true"
	mysqdb, err = sqlx.Connect(dbdrive, dataSourceName)
	if err != nil {
		LogErr("1st Err from dbcnf.mysqInit(): ", err)
	}
	acts := getAvailableAccounts()
	if acts != nil {
		bscElems.Store(mblcach, acts)
	}
	return
}

func DatabaseInit() error {
	return mysqInit()
}

func DatabaseClose() {
	if err := mysqdb.Close(); err != nil {
		LogErr("Err from dbcnf.DatabaseClose(): ", err)
	}
}
