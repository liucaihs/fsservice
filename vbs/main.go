package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"time"
	"vbs/handler"
	"vbs/storage"

	"github.com/gin-gonic/gin"
)

func main() {
	defer func() {
		if recover() != nil {
			var buf [2 << 10]byte
			storage.LogRun("panic.log", "Stack Info is: %s\n", string(buf[:runtime.Stack(buf[:], true)]))
		}
	}()

	if err := storage.DatabaseInit(); err != nil {
		os.Exit(-1)
	}
	defer storage.DatabaseClose()

	gin.SetMode(os.Getenv("GIN_MODE"))
	router := gin.Default()
	router.GET("/1/mobile", handler.ObtainOnlinePhoneNumber)
	router.GET("/vcode/:pkg/:phoneNumber", handler.SyncProvideIdentifyCode)

	srv := &http.Server{
		Addr:        ":8093",
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
