package main

import (
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	pb "rcs/proto"
	"rcs/remote"
)

func main() {
	log.Println("Server start ...")
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", 2800))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterCommandServer(grpcServer, remote.NewRemoteServer())
	go func() {
		grpcServer.Serve(ln)
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit)
	s := <-quit
	log.Println("Got Signal:", s)

	grpcServer.GracefulStop()
	log.Println("Server exit")
}
