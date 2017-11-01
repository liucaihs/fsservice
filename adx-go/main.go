package main

import (
	"adx-go/model"
	"adx-go/remote"
	"log"
	"os"
	"os/signal"
)

func main() {
	if err := model.DatabaseInit(); err != nil {
		panic(err)
		os.Exit(-1)
	}
	log.Println("Server start ...")
	remoteServer := remote.NewRemoteServer()
	remoteServer.StartAll()
	quit := make(chan os.Signal)
	signal.Notify(quit)
	s := <-quit
	log.Println("Got Signal:", s)
	remoteServer.StopAll()
	log.Println("Server exit")
}
