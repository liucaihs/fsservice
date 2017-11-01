package remote

import (
	"adx-go-back/common"
	"adx-go-back/model"
	"log"
	"net/url"
	"strings"
	"time"
)

func notifyChannel(adPut *model.Advertising, activate *model.ActivatedLog, clasp int) {
	log.Println("notifyChannel : active_id", activate.Id, "click_id:", activate.ClickId)
	if clasp > 0 {
		log.Println("notifyChannel", ": 激活扣量不通知渠道！！！！！！。")
		return
	}
	var ordRecd *model.ClickLog

	if activate.ClickId > 0 {
		ordRecd = model.GetClickById(activate.ClickId)
	} else if len(activate.Ifa) > 0 && len(activate.AdputSn) > 0 {
		ordRecd = model.GetClickByIfaSn(activate.Ifa, activate.AdputSn)
	}

	if ordRecd.Id < 1 {
		log.Println("notifyChannel: 找不到点击数据！！！！！！！！！", ordRecd)
		return
	}

	clickRawData, err := common.Json2map([]byte(ordRecd.RawData))
	if err != nil {
		log.Println("notifyChannel Json2map err:", ordRecd.RawData, err)

	}
	if adPut.UseApi > 0 && len(adPut.ChannelUrl) > 0 {
		_, ok := clickRawData["tid"]
		if strings.Contains(adPut.ChannelUrl, "{TID}") && ok {
			adPut.ChannelUrl = strings.Replace(adPut.ChannelUrl, "{TID}", clickRawData["tid"].(string), -1)
		}

		if strings.Contains(adPut.ChannelUrl, "{ifa}") {
			adPut.ChannelUrl = strings.Replace(adPut.ChannelUrl, "{ifa}", ordRecd.Ifa, -1)
		}
	} else {
		_, ok_cu := clickRawData["callback_url"]
		if ok_cu {
			adPut.ChannelUrl = clickRawData["callback_url"].(string)
		}
		_, ok_c := clickRawData["callback"]
		if ok_c {
			adPut.ChannelUrl = clickRawData["callback"].(string)
		}
		if !ok_cu && !ok_c {
			log.Println("notifyChannel: 获取渠道回调地址失败！！！！！！！！！", ordRecd)
			return
		}

	}
	unes_url, err := url.QueryUnescape(adPut.ChannelUrl)
	if err != nil {
		log.Println("notifyChannel ChannelUrl QueryUnescape err:  ！！！！！！！！", adPut.ChannelUrl)
	} else {
		adPut.ChannelUrl = unes_url
	}
	log.Println(adPut.ChannelUrl)
	checkrep_url(adPut.ChannelUrl)

	notifysta := curl_get(adPut.ChannelUrl, activate.Id)

	if notifysta {
		activate.UpNotifyTime(int(time.Now().Unix()))
		log.Println("notifyChannel scusess", adPut.ChannelUrl)
	} else {
		log.Println("notifyChannel fail!!!", adPut.ChannelUrl)
	}

}
