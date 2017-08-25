package remote

import (
	"log"
	"os"
	pb "rcs/proto"
)

type remoteServer struct {
	Logger *log.Logger
}

func (s *remoteServer) ListOnline(o *pb.PageOption, stream pb.Command_ListOnlineServer) error {
	return nil
}

func NewRemoteServer() *remoteServer {
	return &remoteServer{
		Logger: log.New(os.Stderr, "", log.LstdFlags),
	}
}
