package remote

import (
	"adx-go-back/common"
	"adx-go-back/model"
	"log"
	"math/rand"
	"strconv"

	"time"

	"github.com/Shopify/sarama"
)

type mqHandler func(s *remoteServer, message *sarama.ConsumerMessage)

var mqhRouter map[string]mqHandler = map[string]mqHandler{
	"adx_click":    clickHandle,
	"adx_activate": activateHandle,
	"adx_recvifa":  recvifaHandle,
	"adx_short":    shortHandle,
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

func clickHandle(s *remoteServer, message *sarama.ConsumerMessage) {
	inputMap, err := common.Json2map(message.Value)
	if err != nil {
		log.Println("Json2map err:", err)
		return
	}
	log.Println("clickHandle :", string(message.Value), inputMap)
	click := new(model.ClickLog)
	if _, ok := inputMap["s"]; ok {
		click.AdputSn = inputMap["s"].(string)
	}
	if _, ok := inputMap["mac"]; ok {
		click.Mac = inputMap["mac"].(string)
	}
	if _, ok := inputMap["ifa"]; ok {
		click.Ifa = inputMap["ifa"].(string)
	}
	if ipstr, ok := inputMap["ip"]; ok {
		click.Ip = ipstr.(string)
	}
	status, _ := inputMap["status"].(string)
	click.Status, _ = strconv.Atoi(status)
	log.Println("click.Status:", status, inputMap["status"])
	if _, ok := inputMap["FromIp"]; ok {
		click.FromIp = inputMap["FromIp"].(string)
	}

	click.CreateTime = int(time.Now().Unix())
	click.RawData = string(message.Value)
	affected, err := click.Insert()
	if err == nil {
		log.Println("insert:", click, affected)
		log.Println(": 点击入库成功 : click_id : " + strconv.FormatInt(click.Id, 10))

		if _, ok := inputMap["redirect"]; ok {
			if inputMap["redirect"].(int) > 0 {
				click.UpNotifyTime(int(time.Now().Unix()))
				return
			}
		}
		if len(click.AdputSn) > 0 && click.Status == 1 {
			_, adPut := model.GetAdvPuton(click.AdputSn)
			notifyAdvertiser(adPut.AdId, click)
		}

	} else {
		log.Println(": ！！！！！！！点击入库失败！！！！！！！！")
		log.Println("insert: err!!!", click, err)
	}

}

func activateHandle(s *remoteServer, message *sarama.ConsumerMessage) {
	inputMap, err := common.Json2map(message.Value)
	if err != nil {
		log.Println("Json2map err:", err)
		return
	}
	log.Println("activateHandle :", string(message.Value), inputMap)
	activate := new(model.ActivatedLog)
	cid, _ := inputMap["cid"].(string)
	click_id, _ := strconv.ParseInt(cid, 10, 64)
	adput_sn := ""
	in_s, ok := inputMap["s"]
	if ok && len(in_s.(string)) > 0 {
		adput_sn = in_s.(string)
	} else {
		if click_id < 1 {
			return
		}
		click_rec := model.GetClickById(click_id)
		adput_sn = click_rec.AdputSn
	}

	activate.AdputSn = adput_sn
	if _, ok := inputMap["mac"]; ok {
		activate.Mac = inputMap["mac"].(string)
	}
	if _, ok := inputMap["ifa"]; ok {
		activate.Ifa = inputMap["ifa"].(string)
	}
	if click_id < 1 && len(activate.Ifa) > 0 {
		clickRec := model.GetClickByIfaSn(activate.Ifa, activate.AdputSn)
		click_id = clickRec.Id

	}

	activate.Status = 1
	activate.ClickId = click_id

	_, adPut := model.GetAdvPuton(adput_sn)
	var clasp int = 0
	randNum := rand.Intn(100)
	Reduction, err := strconv.ParseFloat(adPut.Reduction, 64)
	dbReduc, err := strconv.Atoi(adPut.Reduction)
	if dbReduc > 0 && float64(randNum) <= Reduction*100 {
		clasp = 1
	}
	activate.Clasp = clasp
	if _, ok := inputMap["FromIp"]; ok {
		activate.FromIp = inputMap["FromIp"].(string)
	}
	activate.RawData = string(message.Value)

	affected, err := activate.Insert()
	if err == nil {

		log.Println("insert activate :", activate, affected)
		log.Println("激活入库成功 : activate_id : " + strconv.FormatInt(activate.Id, 10))
		notifyChannel(adPut, activate, clasp)
	} else {
		log.Println("！！！！！！！激活入库失败！！！！！！！！")
		log.Println("activate insert  err!!!", activate, err)
	}
}

func recvifaHandle(s *remoteServer, message *sarama.ConsumerMessage) {
	inputMap, err := common.Json2map(message.Value)
	if err != nil {
		log.Println("Json2map err:", err)
		return
	}
	log.Println("recvifaHandle :", string(message.Value), inputMap)
	var actIfa *model.ActivatedIfa
	if _, ok := inputMap["ifa"]; ok {
		actIfa.Ifa = inputMap["ifa"].(string)
	}
	if _, ok := inputMap["timestamp"]; ok {
		actIfa.Timestamp = inputMap["timestamp"].(int)
	}
	actIfa.FromIp = inputMap["from_ip"].(string)
	actIfa.RawData = string(message.Value)
	actIfa.CreateTime = int(time.Now().Unix())
	affected, err := actIfa.Insert()
	if err == nil {
		log.Println("insert ActivatedIfa :", actIfa, affected)
		log.Println("ActivatedIfa 激活入库成功 : Id : " + strconv.FormatInt(actIfa.Id, 10))
	} else {
		log.Println("！！！！！！！激活入库失败！！！！！！！！")
		log.Println("ActivatedIfa insert  err!!!", actIfa, err)
	}
}

func shortHandle(s *remoteServer, message *sarama.ConsumerMessage) {
	inputMap, err := common.Json2map(message.Value)
	if err != nil {
		log.Println("Json2map err:", err)
		return
	}
	log.Println("shortHandle :", string(message.Value), inputMap)
	var shortLog *model.ShortLog
	shortLog.CreateTime = int(time.Now().Unix())
	shortLog.RawData = inputMap["code"].(string)
	shortLog.ShortId = inputMap["short_id"].(int)
	shortLog.Ip = inputMap["from_ip"].(string)
	affected, err := shortLog.Insert()
	if err == nil {
		log.Println("insert shortLog :", shortLog, affected)

	} else {
		log.Println("！！！！！！！shortLog入库失败！！！！！！！！")

	}
}
