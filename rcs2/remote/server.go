package remote

import (
	"golang.org/x/net/context"
	"log"
	"os"
	pb "rcs/proto"
)

type remoteServer struct {
	Logger *log.Logger
}

func (s *remoteServer) Signin(context.Context, *pb.ReqSignin) (*pb.RespSignin, error) {
	return nil, nil
}

func (s *remoteServer) CmdTodoVcode(context.Context, *pb.ReportHost) (*pb.None, error) {
	return nil, nil
}

func (s *remoteServer) CmdUploadDevDetail(context.Context, *pb.ReportHost) (*pb.None, error) {
	return nil, nil
}

func NewRemoteServer() *remoteServer {
	return &remoteServer{
		Logger: log.New(os.Stderr, "", log.LstdFlags),
	}
}
