package remote

import (
	"adx-go-back/common"
	"log"
	"os"
	"strings"

	"github.com/Shopify/sarama"
)

var (
	logger = log.New(os.Stderr, "[srama]", log.LstdFlags)
)

type MQ struct {
	ap      sarama.AsyncProducer
	c       sarama.Consumer
	stopped chan bool
}

func newConsumer(brokerList []string) sarama.Consumer {
	sarama.Logger = logger
	consumer, err := sarama.NewConsumer(brokerList, nil)
	if err != nil {
		log.Println("Failed to start sarama consumer:", err)
	}
	return consumer
}

func (mq *MQ) stop() {

	close(mq.stopped)
}

func (mq *MQ) close() {

	if err := mq.c.Close(); err != nil {
		log.Println("Failed to shut down consumer:", err)
	}
}

func (mq *MQ) Run(s *remoteServer, topics []string) {
	if mq.c == nil {
		log.Println("MQ service did not initialize ...")
		return
	}

	for _, topic := range topics {
		partitionList, err := mq.c.Partitions(topic)
		if err != nil {
			logger.Println(topic, "Failed to get the list of partitions: ", err)
		}

		for partition := range partitionList {
			pc, err := mq.c.ConsumePartition(topic, int32(partition), sarama.OffsetNewest) // sarama.OffsetOldest)
			if err != nil {
				logger.Printf(topic, "Failed to start consumer for partition %d: %s\n", partition, err)
			}
			defer pc.AsyncClose()

			go func(pcp sarama.PartitionConsumer, topicName string) {
				for msg := range pcp.Messages() {
					logger.Println(topicName, "message is :", msg)
					logger.Printf("Partition:%d, Offset:%d, Key:%s, Value:%s", msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
					handleMQ(s, topicName, msg)
					logger.Println()
				}
			}(pc, topic)
		}
	}

	defer func() {

		mq.close()
	}()

	for {
		select {

		case <-mq.stopped:

			return
		}
	}
}

func NewMQ() *MQ {
	kafka_server := common.GetEnvDef("KAFKA_ENDPOINTS", "192.168.1.214:9092")
	brokerList := strings.Split(kafka_server, ",")
	log.Printf("kafka brokers: %s", strings.Join(brokerList, ","))
	mq := &MQ{

		c:       newConsumer(brokerList),
		stopped: make(chan bool),
	}

	return mq
}
