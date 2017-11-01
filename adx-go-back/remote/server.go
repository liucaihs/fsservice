package remote

import (
	"log"

	"sync"
)

var Mq *MQ

type remoteServer struct {
	wg *sync.WaitGroup
}

func (s *remoteServer) runMQ() {
	Mq.Run(s, []string{"adx_click", "adx_activate", "adx_recvifa", "adx_short"})
	log.Println("mqServer shutdown ...")
}

func (s *remoteServer) StartAll() {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		s.runMQ()
	}()

}

func (s *remoteServer) StopAll() {
	Mq.stop()
	s.wg.Wait()
}

func NewRemoteServer() *remoteServer {
	r := &remoteServer{
		wg: new(sync.WaitGroup),
	}
	Mq = NewMQ()
	return r
}
