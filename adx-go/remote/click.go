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
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func ClickRecv(c *gin.Context) {
	inputDatas := c.Request.URL.Query()
	transData := urlQueryTrans(inputDatas)
	transData["FromIp"] = c.ClientIP()

	adput_sn := c.DefaultQuery("s", "")
	adput_sn = strings.Replace(adput_sn, " ", "", -1)
	mac := c.DefaultQuery("mac", "")
	mac = strings.Replace(mac, " ", "", -1)
	ifa := c.DefaultQuery("ifa", "")
	ifa = strings.Replace(ifa, " ", "", -1)
	ip := c.DefaultQuery("ip", "")
	ip = strings.Replace(ip, " ", "", -1)

	has, adPut := model.GetAdvPuton(adput_sn)

	if !has {
		showErro(c, "-201")
		return
	}
	if len(ip) > 0 {
		bip := checkIp(ip)
		fmt.Print(bip)
		if !bip {
			showErro(c, "-401")
			return
		}
	}
	if len(mac) == 0 && len(ifa) == 0 {
		showErro(c, "-301")
		return
	}
	clkStatus := model.AdPutonChk(adPut)

	switch clkStatus {
	case 2:
		showErro(c, "-302")
		return
	case 3:
		showErro(c, "-303")
		return
	case 4:
		showErro(c, "-304")
		return
	}
	transData["status"] = strconv.Itoa(clkStatus)
	inputdata, err := json.Marshal(transData)
	if err != nil {
		log.Println(err.Error())
	}
	fmt.Println(inputDatas, transData, inputdata)
	Mq.PutMessage("adx_click", "", string(inputdata))
	if _, ok := transData["goto"]; ok {
		if len(transData["goto"]) > 0 {
			c.Redirect(http.StatusFound, transData["goto"])
		}
	}
	showSuccess(c)
	defer func() {
		if err := recover(); err != nil {
			log.Println(": ！！！！！！！点击入库失败！！！！！！！！", err)
			showErro(c, "")
		}
	}()

}

func clickRedirect(c *gin.Context) {
	inputDatas := c.Request.URL.Query()
	transData := urlQueryTrans(inputDatas)
	transData["FromIp"] = c.ClientIP()

	adput_sn := c.DefaultQuery("s", "")
	adput_sn = strings.Replace(adput_sn, " ", "", -1)
	mac := c.DefaultQuery("mac", "")
	mac = strings.Replace(mac, " ", "", -1)
	ifa := c.DefaultQuery("ifa", "")
	ifa = strings.Replace(ifa, " ", "", -1)
	ip := c.DefaultQuery("ip", "")
	ip = strings.Replace(ip, " ", "", -1)

	has, adPut := model.GetAdvPuton(adput_sn)

	if !has {
		showErro(c, "-201")
		return
	}
	if len(ip) > 0 {
		bip := checkIp(ip)
		fmt.Print(bip)
		if !bip {
			showErro(c, "-401")
			return
		}
	}
	if len(mac) == 0 && len(ifa) == 0 {
		showErro(c, "-301")
		return
	}
	clkStatus := model.AdPutonChk(adPut)
	transData["status"] = strconv.Itoa(clkStatus)
	transData["redirect"] = strconv.Itoa(1)
	inputdata, err := json.Marshal(transData)
	if err != nil {
		log.Println(err.Error())
	}
	fmt.Println(inputDatas, transData, inputdata)
	Mq.PutMessage("adx_click", "", string(inputdata))
	switch clkStatus {
	case 2:
		showErro(c, "-302")
		return
	case 3:
		showErro(c, "-303")
		return
	case 4:
		showErro(c, "-304")
		return
	}

	if strings.Contains(adPut.PutonUrl, "{callback_url}") {
		wequeryStr := "s=" + adput_sn + "&mac=" + transData["mac"] + "&ifa=" + transData["ifa"] + "&cid="
		wecallbackurl := model.GetSetVal("adx_adeff")
		wecallbackurl = wecallbackurl + "?s=" + base64_encode(wequeryStr)
		adPut.PutonUrl = strings.Replace(adPut.PutonUrl, "{callback_url}", wecallbackurl, -1)
	}
	c.Redirect(http.StatusFound, adPut.PutonUrl)
	defer func() {
		if err := recover(); err != nil {
			log.Println(": ！！！！！！！点击入库失败！！！！！！！！", err)
			showErro(c, "")
		}
	}()

}

func shortSpread(c *gin.Context) {
	code := c.Param("code")
	code = strings.Replace(code, " ", "", -1)
	inputDatas := c.Request.URL.Query()
	transData := urlQueryTrans(inputDatas)
	transData["from_ip"] = c.ClientIP()
	transData["code"] = code

	shortCh := model.GetShortBycode(code)
	if shortCh.Id < 1 {
		showNoFound(c, "")
	}
	transData["short_id"] = strconv.Itoa(shortCh.Id)
	inputJson, err := json.Marshal(transData)
	if err != nil {
		log.Println(err.Error())
		showNoFound(c, "")
	}
	Mq.PutMessage("adx_short", "", string(inputJson))
	c.Redirect(http.StatusFound, shortCh.AdvertUrl)
	defer func() {
		if err := recover(); err != nil {
			log.Println(": ！！！！！！！接收失败！！！！！！！！", err)
			showNoFound(c, "")
		}
	}()

}

func notifyAdvertiser(ad_id int, clickLog *model.ClickLog) {
	log.Println(": ！！！！！！！enter notifyAdvertiser！！！！！！！！")
	advm := model.GetById(ad_id)
	fmt.Println(advm)
	if advm == nil {
		log.Println(": 找不到广告信息 ", ad_id)
	}
	wequeryStr := "s=" + clickLog.AdputSn + "&mac=" + clickLog.Mac + "&ifa=" + clickLog.Ifa + "&cid=" + strconv.FormatInt(clickLog.Id, 10)
	wecallbackurl := model.GetSetVal("adx_adeff")
	wecallbackurl = wecallbackurl + "?s=" + base64_encode(wequeryStr)
	fmt.Println(wecallbackurl)

	callurl := advm.OwnerApiurl
	fmt.Println(callurl)
	var rep_strval = ""
	if strings.Contains(callurl, "{ifa_md5}") {
		if len(clickLog.Ifa) > 0 {
			rep_strval = calcmd5(strings.ToUpper(clickLog.Ifa))
		}
		callurl = strings.Replace(callurl, "{ifa_md5}", rep_strval, -1)
	}

	if strings.Contains(callurl, "{ifa}") {
		callurl = strings.Replace(callurl, "{ifa}", clickLog.Ifa, -1)
	}

	if strings.Contains(callurl, "{callback_url}") {
		rep_strval = wecallbackurl
		if advm.Urlencode > 0 {
			rep_strval = url.QueryEscape(wecallbackurl)
		}
		callurl = strings.Replace(callurl, "{callback_url}", rep_strval, -1)
	}

	if strings.Contains(callurl, "{mac_md5}") {
		rep_strval = ""
		if len(clickLog.Mac) > 0 {
			rep_strval = strings.Replace(strings.ToUpper(clickLog.Mac), "-", ":", -1)
			rep_strval = calcmd5(rep_strval)
		}
		callurl = strings.Replace(callurl, "{mac_md5}", rep_strval, -1)
	}

	if strings.Contains(callurl, "{mac}") {
		callurl = strings.Replace(callurl, "{mac}", clickLog.Mac, -1)
	}

	if strings.Contains(callurl, "{time}") {
		callurl = strings.Replace(callurl, "{time}", strconv.Itoa(clickLog.CreateTime), -1)
	}
	if strings.Contains(callurl, "{MAC:}") {
		rep_strval = ""
		if len(clickLog.Mac) > 0 {
			rep_strval = strings.Replace(strings.ToUpper(clickLog.Mac), "-", ":", -1)
			rep_strval = calcmd5(rep_strval)
		}
		callurl = strings.Replace(callurl, "{MAC:}", rep_strval, -1)
	}

	if strings.Contains(callurl, "{ip}") {
		callurl = strings.Replace(callurl, "{ip}", clickLog.Ip, -1)
	}

	fmt.Println(callurl)
	checkrep_url(callurl)
	//	notifysta := curl_get("http://127.0.0.1/test.php", clickLog.Id)
	notifysta := curl_get(callurl, clickLog.Id)
	if notifysta {
		clickLog.UpNotifyTime(int(time.Now().Unix()))
	}
}

func curl_get(url string, cid int64) bool {
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
