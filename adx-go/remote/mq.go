package remote

import (
	"adx-go/common"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Shopify/sarama"
)

var (
	logger = log.New(os.Stderr, "[srama]", log.LstdFlags)
)

type MQ struct {
	ap sarama.AsyncProducer

	stopped    chan bool
	stopped_ap chan bool
}

func newOnoffProducer(brokerList []string) sarama.AsyncProducer {
	sarama.Logger = logger
	config := sarama.NewConfig()
	//	config.Producer.RequiredAcks = sarama.WaitForLocal //Only wait for the leader to ack
	//	config.Producer.Return.Successes = true
	//	config.Producer.Timeout = 5 * time.Second
	//	config.Producer.Compression = sarama.CompressionSnappy
	config.Producer.Flush.Frequency = 500 * time.Millisecond

	producer, err := sarama.NewAsyncProducer(brokerList, config) //sarama.NewAsyncProducer(brokerList, config)
	if err != nil {
		log.Println("Failed to start sarama producer:", err)
	}
	return producer
}

func (mq *MQ) PutMessage(topic string, kstr string, vstr string) {
	if mq.ap == nil {
		log.Println("MQ service did not initialize ...")
		return
	}
	k := sarama.StringEncoder(kstr)
	v := sarama.StringEncoder(vstr)
	log.Println(vstr)
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   k,
		Value: v,
	}

	mq.ap.Input() <- msg

}

func (mq *MQ) stop() {
	close(mq.stopped_ap)
	close(mq.stopped)
}

func (mq *MQ) close() {
	if err := mq.ap.Close(); err != nil {
		log.Println("Failed to shut down async producer", err)
	}

}

func (mq *MQ) ListenPush() {
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
			case <-mq.stopped_ap:
				return
			}
		}
	}()
}

func (mq *MQ) Run(s *remoteServer, topics []string) {
	if mq.ap == nil {
		log.Println("MQ service did not initialize ...")
		return
	}
	mq.ListenPush()
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
		ap:         newOnoffProducer(brokerList),
		stopped:    make(chan bool),
		stopped_ap: make(chan bool),
	}
	return mq
}
