package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"
	"vcode2/handlers"
	"vcode2/storage"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := storage.DatabaseInit(); err != nil {
		os.Exit(-1)
	}

	gin.SetMode(os.Getenv("GIN_MODE"))
	router := gin.Default()
	router.POST("/phones", handlers.CollectPhoneNumber)
	router.GET("/mobiles/:pkgname/:count", handlers.FetchPhoneNumber)
	router.POST("/vcodes", handlers.CollectVerifyCode)
	router.GET("/vcodes/:pkgname", handlers.FetchVerifyCode)
	router.PUT("/setsize", handlers.UpdateSetSize)

	srv := &http.Server{
		Addr:        ":8088",
		Handler:     router,
		IdleTimeout: 2 * time.Minute,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			storage.LogErr("Listen Err: ", err)
		}
	}()
	exit := make(chan os.Signal)
	signal.Notify(exit)
	<-exit
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}
