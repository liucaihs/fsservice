package remote

import (
	"log"

	"github.com/Shopify/sarama"
)

type mqHandler func(s *remoteServer, message *sarama.ConsumerMessage)

var mqhRouter map[string]mqHandler = map[string]mqHandler{
	"vcode":   vcodeHandle,
	"device":  deviceHandle,
	"qualify": qualifyHandle,
}

func handleMQ(s *remoteServer, topic string, message *sarama.ConsumerMessage) {
	handle := mqhRouter[topic]
	if handle == nil {
		log.Println(topic, "handleMQ topic not find ..... :", string(message.Value))
		return
	}
	log.Println(topic, "handleMQ")
	handle(s, message)
}

func vcodeHandle(s *remoteServer, message *sarama.ConsumerMessage) {
	log.Println("vcodeHandle :", string(message.Value))

}

func deviceHandle(s *remoteServer, message *sarama.ConsumerMessage) {
	log.Println("deviceHandle :", string(message.Value))

}

func qualifyHandle(s *remoteServer, message *sarama.ConsumerMessage) {
	log.Println("qualifyHandle :", string(message.Value))

}
