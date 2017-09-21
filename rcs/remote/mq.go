package remote

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/Shopify/sarama"
)

var (
	logger = log.New(os.Stderr, "[srama]", log.LstdFlags)
	topics = "onoff"
)

type MQ struct {
	ap      sarama.AsyncProducer
	c       sarama.Consumer
	stopped chan bool
}

func newOnoffProducer(brokerList []string) sarama.AsyncProducer {
	sarama.Logger = logger
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForLocal //Only wait for the leader to ack
	config.Producer.Return.Successes = true
	config.Producer.Timeout = 5 * time.Second
	config.Producer.Compression = sarama.CompressionSnappy
	config.Producer.Flush.Frequency = 500 * time.Millisecond

	producer, err := sarama.NewAsyncProducer(brokerList, config) //sarama.NewAsyncProducer(brokerList, config)
	if err != nil {
		log.Println("Failed to start sarama producer:", err)
	}
	return producer
}

func newConsumer(brokerList []string) sarama.Consumer {
	sarama.Logger = logger
	consumer, err := sarama.NewConsumer(brokerList, nil)
	if err != nil {
		log.Println("Failed to start sarama consumer:", err)
	}
	return consumer
}

func (mq *MQ) PutMessage(topic string, kstr string, vstr string) { //k sarama.Encoder, v sarama.Encoder) {
	k := sarama.StringEncoder(kstr)
	v := sarama.StringEncoder(vstr)
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   k,
		Value: v,
	}

	mq.ap.Input() <- msg

}

func (mq *MQ) stop() {
	close(mq.stopped)
}

func (mq *MQ) close() {
	if err := mq.ap.Close(); err != nil {
		log.Println("Failed to shut down async producer", err)
	}
	if err := mq.c.Close(); err != nil {
		log.Println("Failed to shut down consumer:", err)
	}
}

func (mq *MQ) Run(s *remoteServer) {
	if mq.c == nil {
		log.Println("MQ service did not initialize ...")
		return
	}
	//必须有这个匿名函数内容
	go func() {
		errors := mq.ap.Errors()
		success := mq.ap.Successes()
		for {
			select {
			case err := <-errors:
				if err != nil {
					log.Printf(err.Error())
				}
			case putres := <-success:
				logger.Printf("partition=%d, offset=%d\n", putres.Partition, putres.Offset)
			}
		}
	}()
	//	pc, err := mq.c.ConsumePartition(topics, 0, sarama.OffsetNewest)
	//	if err != nil {
	//		log.Println("ConsumePartition Failed:", err)
	//		return
	//	}
	partitionList, err := mq.c.Partitions(topics)
	if err != nil {
		logger.Println("Failed to get the list of partitions: ", err)
	}

	for partition := range partitionList {
		pc, err := mq.c.ConsumePartition(topics, int32(partition), sarama.OffsetNewest)
		if err != nil {
			logger.Printf("Failed to start consumer for partition %d: %s\n", partition, err)
		}
		defer pc.AsyncClose()

		go func(sarama.PartitionConsumer) {
			for msg := range pc.Messages() {
				logger.Println("message is :", msg)
				logger.Printf("Partition:%d, Offset:%d, Key:%s, Value:%s", msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
				handleMQ(s, msg)
				logger.Println()
			}
		}(pc)
	}

	defer func() {
		//		if err := pc.Close(); err != nil {
		//			log.Println("Failed to Close pc:", err)
		//		}
		mq.close()
	}()

	for {
		select {
		//		case message := <-pc.Messages():
		//			log.Println("MQ Consume one:", string(message.Value))
		//			handleMQ(s, message)
		case <-mq.stopped:
			return
		}
	}
}

func NewMQ() *MQ {
	brokerList := strings.Split("192.168.1.214:9092", ",")
	log.Printf("kafka brokers: %s", strings.Join(brokerList, ","))
	return &MQ{
		ap:      newOnoffProducer(brokerList),
		c:       newConsumer(brokerList),
		stopped: make(chan bool),
	}
}
