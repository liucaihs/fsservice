package remote

import (
	"adx-go/model"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"net/http"
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AdeffRecv(c *gin.Context) {
	inputStr := c.DefaultQuery("s", "")
	queryStr := base64_decode(inputStr)
	fmt.Println(queryStr)
	inputData, _ := url.ParseQuery(queryStr)

	transData := urlQueryTrans(inputData)
	cid, ok := transData["cid"]
	if ok {
		transData["cid"] = cid
	} else {
		transData["cid"] = ""
	}
	transData["FromIp"] = c.ClientIP()
	inputJson, err := json.Marshal(transData)
	if err != nil {
		log.Println(err.Error())
	}

	log.Println("AdeffRecv input : ", cid, inputData)
	adput_sn := inputData["s"][0]
	mac := inputData["mac"][0]
	ifa := inputData["ifa"][0]
	fmt.Println(mac, ifa)
	if len(mac) == 0 && len(ifa) == 0 {
		showErro(c, "-301")
		return
	}
	has, adPut := model.GetAdvPuton(adput_sn)

	if !has || adPut.Id < 1 {
		showErro(c, "-201")
		return
	}
	Mq.PutMessage("adx_activate", "", string(inputJson))
	showSuccess(c)
	defer func() {
		if err := recover(); err != nil {
			log.Println(": ！！！！！！！激活入库失败！！！！！！！！", err)
			showErro(c, "")
		}
	}()

}

func adeffRecvFixed(c *gin.Context) {
	inputDatas := c.Request.URL.Query()
	transData := urlQueryTrans(inputDatas)

	cid, ok := transData["cid"]
	if ok {
		transData["cid"] = cid
	} else {
		transData["cid"] = ""
	}

	adput_sn := ""
	mac := transData["mac"]
	ifa := transData["ifa"]
	log.Println("AdeffRecv input : ", cid, transData, mac, ifa)
	var clickRecord *model.ClickLog
	if len(cid) > 0 {
		cid64, _ := strconv.ParseInt(cid, 10, 64)
		clickRecord = model.GetClickById(cid64)
	} else if len(ifa) > 0 {
		clickRecord = model.GetClickByIfa(ifa)
		transData["cid"] = strconv.FormatInt(clickRecord.Id, 10)
	}
	adput_sn = clickRecord.AdputSn
	has, adPut := model.GetAdvPuton(adput_sn)

	if !has || adPut.Id < 1 {
		showErro(c, "-201")
		return
	}
	transData["FromIp"] = c.ClientIP()
	inputJson, err := json.Marshal(transData)
	if err != nil {
		log.Println(err.Error())
	}
	Mq.PutMessage("adx_activate", "", string(inputJson))
	showSuccess(c)
	defer func() {
		if err := recover(); err != nil {
			log.Println(": ！！！！！！！激活入库失败！！！！！！！！", err)
			showErro(c, "")
		}
	}()

}

func recvIfa(c *gin.Context) {
	inputDatas := c.Request.URL.Query()
	transData := urlQueryTrans(inputDatas)
	transData["from_ip"] = c.ClientIP()
	inputJson, err := json.Marshal(transData)
	if err != nil {
		log.Println(err.Error())
	}
	Mq.PutMessage("adx_recvifa", "", string(inputJson))
	showSuccess(c)
	defer func() {
		if err := recover(); err != nil {
			log.Println(": ！！！！！！！接收失败！！！！！！！！", err)
			showErro(c, "")
		}
	}()

}

func notifyChannel(adPut *model.Advertising, inputData map[string][]string, active_id int64, clasp int) {
	click_id, _ := strconv.ParseInt(inputData["cid"][0], 10, 64)
	fmt.Println("input : active_id : ", active_id, "click_id : ", click_id)
	if clasp > 0 {
		fmt.Println(": 激活扣量不通知渠道！！！！！！。")
	}

}

func curl_get_ac(url string, cid int64) bool {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get(url)

	if err != nil {
		fmt.Println("error:", err)
		return false
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	//	fmt.Println(string(body))
	if len(body) > 200 {
		body = body[0:200]
	}
	fmt.Println(len(string(body)))
	fmt.Println(string(body))
	status := resp.StatusCode
	fmt.Println(status)
	if status >= 400 {
		log.Println("notifyAdvertiser fail!!!", status, url)
		return false
	} else {
		log.Println("notifyAdvertiser scusess", status, url)
		return true
	}

}
