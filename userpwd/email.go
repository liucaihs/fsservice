package main

import (
	"os"
	"strconv"
	"sync"
	"time"

	"gopkg.in/gomail.v2"
)

const (
	zeroPoint = "perdaytmstp"
	usednum   = "perdaycount"
	mblcach   = "csvinfos"
)

var bscElems sync.Map

func Alert() {
	go updateZeroPoint()
	var count int
	for {
		nwusd, ok := bscElems.Load(usednum)
		if ok && nwusd.(int) > 0 {
			query := "select count(*) from view_unused"
			if err := mysqdb.Get(&count, query); err != nil {
				LogErr("1st Err from email.Alert(): ", err)
				continue
			}
			if count <= 1000 {
				notifyITIL(count)
				bscElems.Store(usednum, nwusd.(int)-1)
			}
		}
		time.Sleep(5 * time.Minute)
	}
}

func notifyITIL(nwcont int) {
	numStr := strconv.Itoa(nwcont) + " 个! "
	m := gomail.NewMessage()
	m.SetHeader("From", "migou2306@126.com")
	m.SetHeader("To", os.Getenv("ALERT_EMAIL"))
	m.SetHeader("Subject", "紧急通知")
	m.SetBody("text/plain", "亲，当前可用的用户名个数仅剩"+numStr+numStr+numStr)
	d := gomail.NewDialer("smtp.126.com", 25, "migou2306@126.com", "ruiyang353387")
	if err := d.DialAndSend(m); err != nil {
		LogErr("1st Err from email.notifyITIL(): ", err)
	}
}

func updateZeroPoint() {
	for {
		olderZp, ok := bscElems.Load(zeroPoint)
		if !ok || time.Now().Unix()-olderZp.(int64) > 24*3600 {
			zeroTM, _ := time.ParseInLocation("2006-01-02", time.Now().Format("2006-01-02"), time.Local)
			LogRun("email_raw_datas.log", "rawdata: %v %T\n", zeroTM, zeroTM)
			LogRun("email_raw_datas.log", "rawdata: %v %T\n", zeroTM.Unix(), zeroTM.Unix())
			bscElems.Store(zeroPoint, zeroTM.Unix())
			bscElems.Store(usednum, 3)
		}
		time.Sleep(1 * time.Hour)
	}
}
