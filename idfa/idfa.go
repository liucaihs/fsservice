package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var db *sqlx.DB

func getEnvDef(name string, def_str string) string {
	value := os.Getenv(name)
	if len(value) > 0 {
		return value
	}
	return def_str
}

type Idfa struct {
	Ifa         string `json:"ifa" form:"ifa" binding:"required"`
	Mac         string `json:"mac" form:"mac"`
	Ip          string `json:"ip" form:"ip"`
	Create_time int64  `json:"create_time" form:"create_time"`
}

func initDb() {
	db_con := getEnvDef("DB_CONNECTION", "mysql")
	dbuser := getEnvDef("DB_USERNAME", "root")
	dbpwd := getEnvDef("DB_PASSWORD", "secretpw")
	dbhost := getEnvDef("DB_HOST", "192.168.1.211")
	dbport := getEnvDef("DB_PORT", "3306")
	dbname := getEnvDef("DB_DATABASE", "idfa")

	//	var dsn = "root:secretpw@tcp(192.168.1.211:3306)/idfa?charset=utf8"
	var err error
	var dataSourceName string = dbuser + ":" + dbpwd + "@tcp(" + dbhost + ":" + dbport + ")/" + dbname + "?parseTime=true"
	fmt.Println(dataSourceName)
	db, err = sqlx.Connect(db_con, dataSourceName)

	checkErr(err)

	db.SetMaxIdleConns(1000)
	db.SetMaxOpenConns(1000)

}

func main() {
	initDb()
	defer db.Close()

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default() //获得路由实例

	//添加中间件
	//	router.Use(Middleware)
	//注册接口
	router.GET("/", GetHandler)
	router.POST("/collect", PostHandler)

	//监听端口
	srv := &http.Server{
		Addr:        ":8005",
		Handler:     router,
		IdleTimeout: 1 * time.Minute,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			panic(err)
		}
	}()
	log.Println("Server start at 8005")
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, os.Kill)
	s := <-quit
	log.Printf("Shutdown Server ...(signal is %v)\n", s)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		panic(err)
		return
	}

}

func checkErr(errMasg error) {
	if errMasg != nil {
		panic(errMasg)
	}
}
func Middleware(c *gin.Context) {
	fmt.Println("this is a middleware!")
}

func GetHandler(c *gin.Context) {
	c.Data(http.StatusOK, "text/plain", []byte(fmt.Sprintf("get success! %s\n", " ")))
	return
}

type Inputdata struct {
	Ifa string `json:"ifa" form:"ifa" binding:"required"`
	Mac string `json:"mac" form:"mac"`
	Ip  string `json:"ip"  form:"ip"`
}

func PostHandler(c *gin.Context) {
	type JsonHolder struct {
		Status string `json:"status"`
		Msg    string `json:"msg"`
	}
	var inData Inputdata
	if err := c.Bind(&inData); err != nil {
		fmt.Println("Err from  input json data ")
		return
	}

	defer c.Request.Body.Close()

	ifa := inData.Ifa
	ifa = strings.Replace(ifa, " ", "", -1)
	holder := JsonHolder{Status: "fail", Msg: "insert fail"}
	if strings.Compare(ifa, "") != 0 {
		var dvcinf Idfa = Idfa{Ifa: ifa, Create_time: time.Now().Unix()}
		lastId, _ := dvcinf.InsertADeviceInfo()
		if lastId > 0 {
			holder = JsonHolder{Status: "ok", Msg: "success"}
		}
	} else {

		holder = JsonHolder{Status: "fail", Msg: "ifa empty"}
	}

	c.JSON(http.StatusOK, holder)
	return
}

func (dvif *Idfa) InsertADeviceInfo() (int64, error) {

	tx, err := db.Beginx()
	if err != nil {
		fmt.Println("Beginx error:", err)
		panic(err)
	}
	sql := "insert into tb_ifa_log (ifa, mac ,ip , create_time) values (:ifa, :mac , :ip , :create_time)"
	fmt.Print(dvif)

	insResult, err := tx.NamedExec(sql, dvif)
	if err != nil {
		tx.Rollback()
		fmt.Println("Beginx Rollback:", err)
		panic(err)
	}
	tx.Commit()
	return insResult.LastInsertId()
}
