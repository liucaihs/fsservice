package model

import (
	"log"
)

type ClickLog struct {
	Id             int64  `xorm:"pk autoincr BIGINT(20)"`
	AdputSn        string `xorm:"not null index VARCHAR(255)"`
	Mac            string `xorm:"not null VARCHAR(255)"`
	Ifa            string `xorm:"VARCHAR(255)"`
	Ip             string `xorm:"VARCHAR(255)"`
	Status         int    `xorm:"default 1 TINYINT(2)"`
	FromIp         string `xorm:"not null VARCHAR(255)"`
	RawData        string `xorm:"not null index TEXT"`
	CreateTime     int    `xorm:"index INT(11)"`
	CallAdvertiser int    `xorm:"INT(11)"`
}

func (c *ClickLog) Insert() (int64, error) {
	return engine.Insert(c)
}

func (c *ClickLog) UpNotifyTime(calltime int) {
	c.CallAdvertiser = calltime
	affected, err := engine.ID(c.Id).Update(&ClickLog{CallAdvertiser: calltime})
	if affected < 1 || err != nil {
		log.Println("update ClickLog CallAdvertiser fail !!!!!!!")
	}
}

func GetClickById(id int64) *ClickLog {
	click := &ClickLog{}
	engine.Id(id).Get(click)
	return click
}

func GetClickByIfa(ifa string) *ClickLog {
	click := new(ClickLog)
	engine.Where("ifa = ?", ifa).Get(click)
	return click
}
