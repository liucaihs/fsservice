package main

import (
	"context"
	"device-backend/handler"
	"device-backend/storage"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := storage.DatabaseInit(); err != nil {
		os.Exit(-1)
	}
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.GET("/deviceinfo/:applicationame", handler.DeviceInfoDisplay)

	srv := &http.Server{
		Addr:        ":8091",
		Handler:     router,
		IdleTimeout: 2 * time.Minute,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			storage.LogErr("Listen Err: ", err)
		}
	}()
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, os.Kill)
	s := <-quit
	storage.LogRun("tmpdata.log", "Shutdown Server ...(signal is %v)\n", s)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		storage.LogErr("Server shutdown Err: ", err)
		return
	}
}
