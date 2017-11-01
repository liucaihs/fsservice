package model

import (
	"log"
	"strconv"
	"time"
)

type ActivatedLog struct {
	Id          int64  `xorm:"pk autoincr BIGINT(20)"`
	AdputSn     string `xorm:"not null VARCHAR(255)"`
	Mac         string `xorm:"not null VARCHAR(255)"`
	Ifa         string `xorm:"VARCHAR(255)"`
	Status      int    `xorm:"default 0 TINYINT(1)"`
	ClickId     int64  `xorm:"BIGINT(20)"`
	CreateTime  int    `xorm:"INT(11)"`
	NotifyMedia int    `xorm:"INT(11)"`
	FromIp      string `xorm:"not null VARCHAR(255)"`
	RawData     string `xorm:"not null TEXT"`
	Clasp       int    `xorm:"default 0 TINYINT(2)"`
}

func (c *ActivatedLog) Insert() (int64, error) {
	return engine.Insert(c)
}

func (c *ActivatedLog) UpNotifyTime(calltime int) {
	c.NotifyMedia = calltime
	affected, err := engine.ID(c.Id).Update(&ActivatedLog{NotifyMedia: calltime})
	if affected < 1 || err != nil {
		log.Println("update ActivatedLog NotifyMedia fail !!!!!!!")
	}
}
func getTodayLim(adput_sn string) int {
	today := time.Now().Format("2017-01-02")

	sql := "select count(distinct raw_data) as ct from tb_activated_log where adput_sn = '" + adput_sn + "' and clasp = 0 and FROM_UNIXTIME(create_time, '%Y-%m-%d') ='" + today + "'"
	results, err := engine.QueryString(sql)
	if err != nil {
		log.Println(err.Error())
	}
	ct, err := strconv.Atoi(results[0]["ct"])
	return ct
}
