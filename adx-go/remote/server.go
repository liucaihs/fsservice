package remote

import (
	"log"

	"sync"
)

var Mq *MQ

type remoteServer struct {
	httpgs *Httpgs
	wg     *sync.WaitGroup
}

func (s *remoteServer) runMQ() {
	Mq.Run(s, []string{})
	log.Println("mqServer shutdown ...")
}

func (s *remoteServer) runHttpgs() {
	s.httpgs.Run()
	log.Println("Http grpc Server shutdown ...")
}

func (s *remoteServer) StartAll() {
	s.wg.Add(2)
	go func() {
		defer s.wg.Done()
		s.runMQ()
	}()

	go func() {
		defer s.wg.Done()
		s.runHttpgs()
	}()
}

func (s *remoteServer) StopAll() {
	Mq.stop()
	s.httpgs.Stop()
	s.wg.Wait()
}

func NewRemoteServer() *remoteServer {
	r := &remoteServer{
		httpgs: NewHttpgs(),
		wg:     new(sync.WaitGroup),
	}
	Mq = NewMQ()
	return r
}
