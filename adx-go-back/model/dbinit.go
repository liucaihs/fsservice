package model

import (
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
)

var engine *xorm.Engine

func getEnvDef(name string, def_str string) string {
	value := os.Getenv(name)
	if len(value) > 0 {
		return value
	}
	return def_str
}

func DatabaseInit() (err error) {
	db_con := getEnvDef("DB_CONNECTION", "mysql")
	dbuser := getEnvDef("DB_USERNAME", "root")
	dbpwd := getEnvDef("DB_PASSWORD", "secretpw")
	dbhost := getEnvDef("DB_HOST", "192.168.1.211")
	dbport := getEnvDef("DB_PORT", "3306")
	dbname := getEnvDef("DB_DATABASE", "adp")

	var dataSourceName string = dbuser + ":" + dbpwd + "@tcp(" + dbhost + ":" + dbport + ")/" + dbname + "?charset=utf8"
	engine, err = xorm.NewEngine(db_con, dataSourceName) //"root:123@/adp2?charset=utf8")
	if err != nil {
		return err
	}
	tbMapper := core.NewPrefixMapper(core.SnakeMapper{}, "tb_")
	engine.SetTableMapper(tbMapper)
	engine.ShowSQL(true)
	//	err = engine.Sync2(new(Advertising))
	if err != nil {
		return err
	}

	return err
}
