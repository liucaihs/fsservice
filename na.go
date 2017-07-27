package main

import "gopkg.in/gin-gonic/gin.v1"
import "net/http"
import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func query_row(dsn string, strsql string) string {
	db, err := sql.Open("mysql", dsn)
	var ret string
	if err != nil {
		log.Println(err)
		return ret
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Println(err)
		return ret
	}

	row := db.QueryRow(strsql)
	row.Scan(&ret)

	return ret
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	var dsn = "na_user:na_pwd001@tcp(rm-wz9ljy2y85c0g191xo.mysql.rds.aliyuncs.com:3306)/na?charset=utf8"
	router := gin.Default()

	v1 := router.Group("/api/v1")
	{
		v1.GET("/nick", func(c *gin.Context) {
			nick := query_row(dsn, "select nick from tb_nick order by RAND() limit 1")
			c.JSON(http.StatusOK, gin.H{"nick": nick})
		})
		v1.GET("/avatar", func(c *gin.Context) {
			avatar := query_row(dsn, "select url from tb_avatar order by RAND() limit 1")
			c.JSON(http.StatusOK, gin.H{"avatar": avatar})
		})

		v1.GET("/avatar/:x/:y", func(c *gin.Context) {

			xval := c.Param("x")
			yval := c.Param("y")
			var condition string = " where 1=1"
			if len(xval) > 0 {
				condition = condition + " and x>=" + xval
			}
			if len(yval) > 0 {
				condition = condition + " and y>=" + yval
			}
			sqlStr := "select url from tb_avatar" + condition + " order by RAND() limit 1"
			avatar := query_row(dsn, sqlStr)
			c.JSON(http.StatusOK, gin.H{"avatar": avatar})
		})
	}

	router.Run(":8500")
}
