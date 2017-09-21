package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := DatabaseInit(); err != nil {
		os.Exit(-1)
	}
	defer DatabaseClose()
	go Alert()

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.GET("/account", SendAccount)
	router.POST("/result", UploadResult)
	srv := &http.Server{
		Addr:        ":10001",
		Handler:     router,
		IdleTimeout: 2 * time.Minute,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			LogErr("server Listen Err: ", err)
		}
	}()
	quit := make(chan os.Signal)
	signal.Notify(quit)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 9*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}
