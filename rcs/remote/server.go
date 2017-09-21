package remote

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	pb "rcs/proto"
	"sync"
	"time"

	xcontext "golang.org/x/net/context"
	"google.golang.org/grpc"
)

type remoteServer struct {
	model  *Model
	mq     *MQ
	gs     *grpc.Server
	ws     *http.Server
	hub    *Hub
	wg     *sync.WaitGroup
	Logger *log.Logger
}

func (s *remoteServer) ListOnline(ctx xcontext.Context, in *pb.GrpcRequest) (*pb.GrpcReply, error) {
	log.Println(in)
	s.model.GetRowNums("")
	return handleGrpc(s, in) //&pb.GrpcReply{TotalResults: 10, HasNext: true, Data: `[{"aaaaa":"bbbbbbb"}]}`}, nil
}

func (s *remoteServer) runMQ() {
	s.mq.Run(s)
	log.Println("mqServer shutdown ...")
}

func (s *remoteServer) runServeGrpc() {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", 2800))
	if err != nil {
		log.Println("grpcServer failed to listen: ", err)
		return
	}

	pb.RegisterCommandServer(s.gs, s)
	s.gs.Serve(ln)
	log.Println("grpcServer shutdown ...")
}

func (s *remoteServer) runServeWs() {
	mux := http.NewServeMux()
	mux.HandleFunc("/rc", func(w http.ResponseWriter, r *http.Request) {
		s.hub.serveWs(w, r)
	})
	s.ws.Handler = mux
	if err := s.ws.ListenAndServe(); err != nil {
		log.Println("wsServer ListenAndServe:", err)
	}
	log.Println("wsServer shutdown ...")
}

func (s *remoteServer) stopServeWs() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.ws.Shutdown(ctx); err != nil {
		log.Println("stopServeWs:", err)
	}
}

func (s *remoteServer) runHub() {
	s.hub.runLoop()
	log.Println("hub runLoop stopped ...")
}

func (s *remoteServer) StartAll() {
	s.wg.Add(4)
	go func() {
		defer s.wg.Done()
		s.runMQ()
	}()
	go func() {
		defer s.wg.Done()
		s.runServeGrpc()
	}()
	go func() {
		defer s.wg.Done()
		s.runServeWs()
	}()
	go func() {
		defer s.wg.Done()
		s.runHub()
	}()
}

func (s *remoteServer) StopAll() {
	s.mq.stop()
	s.gs.GracefulStop()
	s.stopServeWs()
	s.model.CloseConn()
	s.hub.stop()
	s.wg.Wait()
}

func NewRemoteServer() *remoteServer {
	r := &remoteServer{
		model:  GetModel(),
		mq:     NewMQ(),
		gs:     grpc.NewServer(),
		ws:     &http.Server{Addr: ":9090", Handler: nil},
		hub:    newHub(),
		wg:     new(sync.WaitGroup),
		Logger: log.New(os.Stderr, "", log.LstdFlags),
	}
	r.hub.Srv = r
	return r
}
