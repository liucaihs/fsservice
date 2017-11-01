package remote

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	address = "localhost:2800"
)

type Httpgs struct {
	srv     *http.Server
	stopped chan bool
}

func (hpgs *Httpgs) Run() {
	go func() {
		if err := hpgs.srv.ListenAndServe(); err != nil {
			log.Println("server Listen Err: ", err)
		}
	}()
	defer hpgs.Close()
	for {
		select {
		case <-hpgs.stopped:
			return
		}
	}
}

func (hpgs *Httpgs) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), 9*time.Second)
	defer cancel()
	hpgs.srv.Shutdown(ctx)
}

func (hpgs *Httpgs) Stop() {
	close(hpgs.stopped)
}

func NewHttpgs() *Httpgs {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.GET("/adx/click", ClickRecv)
	router.GET("/adx/clickRedirect", clickRedirect)
	router.GET("/adx/spread/:code", shortSpread)

	router.GET("/adx/activate", AdeffRecv)
	router.GET("/adx/activity", adeffRecvFixed)
	router.GET("/adx/recv_ifa", recvIfa)

	router.POST("/onlineList", OnlineList)
	srv := &http.Server{
		Addr:        ":8080",
		Handler:     router,
		IdleTimeout: 5 * time.Minute,
	}

	httpgs := &Httpgs{
		srv:     srv,
		stopped: make(chan bool),
	}

	return httpgs
}

type Inputdata struct {
	PageSize int32  `json:"pageSize"  binding:"required"`
	PageNo   int32  `json:"pageNo"  binding:"required"`
	Extend   string `json:"extend" form:"extend"`
}

func OnlineList(c *gin.Context) {

	var inData Inputdata
	if err := c.Bind(&inData); err != nil {
		log.Println("Err from  input json data ", inData)
		return
	}
	log.Println("input json data:", inData.PageNo, inData.PageSize, inData.Extend)

	c.JSON(http.StatusOK, inData)
}
