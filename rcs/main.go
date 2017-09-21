package main

import (
	"log"
	"os"
	"os/signal"
	"rcs/remote"
)

func main() {
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
