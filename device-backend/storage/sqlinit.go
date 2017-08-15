package storage

import (
	"os"

	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB
var cache *redis.Client

func initRedis() {
	cache = redis.NewClient(&redis.Options{
		Addr:       "redis_server:6379",
		Password:   "",
		DB:         0,
		PoolSize:   999,
		MaxRetries: 3,
	})

}

func initMysql() (err error) {
	host := os.Getenv("MYSQL_SERVER_HOST")
	port := os.Getenv("MYSQL_SERVER_PORT")
	username := os.Getenv("MYSQL_LOGIN_USER")
	password := os.Getenv("MYSQL_USER_PASSWORD")

	var dataSourceName string = username + ":" + password + "@tcp(" + host + ":" + port + ")/device?parseTime=true"
	db, err = sqlx.Connect("mysql", dataSourceName)
	if err != nil {
		LogErr("Err from storage.sqlinit.initMysql(): ", err)
		return err
	}
	db.SetMaxIdleConns(1000)
	db.SetMaxOpenConns(1000)
	return nil
}

func DatabaseInit() (err error) {
	if err = initMysql(); err != nil {
		return err
	}
	initRedis()
	return nil
}
