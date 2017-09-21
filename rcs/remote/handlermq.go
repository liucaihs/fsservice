package remote

import (
	"log"

	"github.com/Shopify/sarama"
)

func handleMQ(s *remoteServer, message *sarama.ConsumerMessage) {
	log.Println(" handleMQ :", string(message.Value))
	inputData, err := Json2map(message.Value)
	if err != nil {
		log.Println("handleMQ Unmarsha1 err:", err)
	}
	log.Println(inputData)
	cliExt := s.hub.Srv.model.GetClintInfo(inputData["phone"].(string), "")
	log.Println("cliExt info:", cliExt)
	if sendid, ok := cliExt["id"]; ok {
		selClient, cok := s.hub.clients[uint32(sendid.(int))]
		log.Println("selClient info:", selClient)
		if cok {
			var sendmsg BaseMessage
			sendmsg.Cmd = "Vcode"
			sendmsg.Data = message.Value
			selClient.sendMessage(&sendmsg)

		}
	}
}
