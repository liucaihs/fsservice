package remote

import (
	"adx-go-back/model"

	"encoding/json"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func notifyAdvertiser(ad_id int, clickLog *model.ClickLog) {
	if ad_id == 44 {
		notifyJdAdvertiser(ad_id, clickLog)
		return
	}
	log.Println(": ！！！！！！！enter notifyAdvertiser！！！！！！！！")
	advm := model.GetById(ad_id)
	log.Println(advm)
	if advm == nil {
		log.Println(": 找不到广告信息 ", ad_id)
	}
	wequeryStr := "s=" + clickLog.AdputSn + "&mac=" + clickLog.Mac + "&ifa=" + clickLog.Ifa + "&cid=" + strconv.FormatInt(clickLog.Id, 10)
	wecallbackurl := model.GetSetVal("adx_adeff")
	wecallbackurl = wecallbackurl + "?s=" + base64_encode(wequeryStr)
	log.Println(wecallbackurl)

	callurl := advm.OwnerApiurl
	log.Println(callurl)
	var rep_strval = ""
	if strings.Contains(callurl, "{ifa_md5}") {
		if len(clickLog.Ifa) > 0 {
			rep_strval = calcmd5(strings.ToUpper(clickLog.Ifa))
		}
		callurl = strings.Replace(callurl, "{ifa_md5}", rep_strval, -1)
	}
	if strings.Contains(callurl, "{clickid}") {
		callurl = strings.Replace(callurl, "{clickid}", url.QueryEscape(wecallbackurl), -1)
	}
	if strings.Contains(callurl, "{cid}") {
		callurl = strings.Replace(callurl, "{cid}", strconv.FormatInt(clickLog.Id, 10), -1)
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

	log.Println(callurl)
	checkrep_url(callurl)
	//	notifysta := curl_get("http://127.0.0.1/test.php", clickLog.Id)
	notifysta := curl_get(callurl, clickLog.Id)
	if notifysta {
		log.Println("notifyAdvertiser scusess", callurl)
		clickLog.UpNotifyTime(int(time.Now().Unix()))
	} else {
		log.Println("notifyAdvertiser fail!!!", callurl)
	}
}

func notifyJdAdvertiser(ad_id int, clickLog *model.ClickLog) {
	log.Println(": ！！！！！！！enter notifyJdAdvertiser！！！！！！！！")
	advm := model.GetById(ad_id)
	log.Println(advm)
	if advm == nil {
		log.Println(": 找不到广告信息 ", ad_id)
		return
	}
	callurl := "http://advch.o2o.jd.com/adchannel/clickV2"
	appid := "jd-o2o"
	source := "mitu-171017rm03"

	wequeryStr := "s=" + clickLog.AdputSn + "&mac=" + clickLog.Mac + "&ifa=" + clickLog.Ifa + "&cid=" + strconv.FormatInt(clickLog.Id, 10)
	wecallbackurl := model.GetSetVal("adx_adeff")
	wecallbackurl = wecallbackurl + "?s=" + base64_encode(wequeryStr)
	log.Println(wecallbackurl)
	post_data := make(map[string]string)
	post_data["appid"] = appid
	post_data["source"] = source
	post_data["idfa"] = clickLog.Ifa
	post_data["ip"] = clickLog.Ip
	post_data["eventtime"] = time.Now().Format("2006-01-02 15:04:05")
	post_data["callbackUrl"] = wecallbackurl
	reqdata, err := json.Marshal(post_data)
	if err != nil {
		log.Println(err.Error())
	}
	req_res := HttpPost2(callurl, "clickParam="+string(reqdata))
	if req_res {
		clickLog.UpNotifyTime(int(time.Now().Unix()))
	}
}
