package remote

import (
	"errors"
	"log"
	"math"
	pb "rcs/proto"
)

type GrpcHandler func(s *remoteServer, in *pb.GrpcRequest, out *pb.GrpcReply) error

var grpcRouter map[string]GrpcHandler = map[string]GrpcHandler{
	"online_list": onlineList,
}

func handleGrpc(s *remoteServer, in *pb.GrpcRequest) (*pb.GrpcReply, error) {
	var grpcReply pb.GrpcReply
	var err error = nil

	if len(in.Action) < 1 {
		err = errors.New("handleGrpc input Action err")
		return &grpcReply, err
	}
	handle := grpcRouter[in.Action]
	if handle == nil {
		err = errors.New("handleGrpc input Action no find")
		return &grpcReply, err
	}
	err = handle(s, in, &grpcReply)
	return &grpcReply, err
}
func checkPage(in *pb.GrpcRequest) error {
	if in.PageNo < 1 {
		log.Println("checkPage PageNo =", in.PageNo)
		in.PageNo = 1
	}
	if in.PageSize < 1 {
		return errors.New("PageSize value error")
	}
	return nil
}

func onlineList(s *remoteServer, in *pb.GrpcRequest, out *pb.GrpcReply) error {

	log.Println("handleGrpc input Action no find.......")
	err := checkPage(in)
	if err != nil {
		return err
	}
	datajson, counts := s.model.GetPageRec()
	out.Data = datajson
	out.TotalResults = counts
	pages := math.Ceil(float64(counts / in.PageSize))
	if float64(in.PageNo) < pages {
		out.HasNext = true
	}
	return nil
}
