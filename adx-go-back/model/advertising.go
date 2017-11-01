package model

import (
	"time"
)

type Advertising struct {
	Id            int    `xorm:"not null pk autoincr INT(10)"`
	AdputSn       string `xorm:"not null VARCHAR(55)"`
	AdputName     string `xorm:"not null default '' VARCHAR(55)"`
	AdId          int    `xorm:"not null INT(11)"`
	AdName        string `xorm:"VARCHAR(255)"`
	MediaId       int    `xorm:"not null INT(11)"`
	MediaName     string `xorm:"not null VARCHAR(85)"`
	PayType       int    `xorm:"not null default 0 TINYINT(3)"`
	Price         int    `xorm:"not null INT(10)"`
	PutonUrl      string `xorm:"VARCHAR(155)"`
	ChannelUrl    string `xorm:"VARCHAR(155)"`
	UseApi        int    `xorm:"not null default 0 TINYINT(2)"`
	StartTime     int    `xorm:"INT(11)"`
	IntoPrice     int    `xorm:"not null default 0 INT(10)"`
	EndTime       int    `xorm:"INT(11)"`
	PutonPlatform int    `xorm:"not null default 1 TINYINT(1)"`
	Limited       int    `xorm:"not null default 0 INT(11)"`
	Remark        string `xorm:"default '' VARCHAR(255)"`
	Status        int    `xorm:"default 1 TINYINT(1)"`
	CreateTime    int    `xorm:"INT(11)"`
	UpdateTime    int    `xorm:"INT(11)"`
	Reduction     string `xorm:"default 0.0000 DECIMAL(6,4)"`
}

func GetAdvPuton(adput_sn string) (bool, *Advertising) {
	adPut := new(Advertising)
	has, err := engine.Where("adput_sn=?", adput_sn).Get(adPut)
	if err != nil {
		has = false
	}
	return has, adPut
}

func AdPutonChk(adPut *Advertising) int {
	var clkStatus = 1
	if adPut.Status == 2 {
		clkStatus = 2
	}
	curTime := int(time.Now().Unix())
	if adPut.StartTime > curTime || adPut.EndTime < curTime {
		clkStatus = 3
	}
	if adPut.Limited > 0 {
		active_nums := getTodayLim(adPut.AdputSn)
		if active_nums > adPut.Limited {
			clkStatus = 4
		}
	}

	return clkStatus
}
