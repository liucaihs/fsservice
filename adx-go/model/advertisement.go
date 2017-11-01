package model

import (
	"log"
)

type Advertisement struct {
	Id            int    `xorm:"not null pk autoincr INT(10)"`
	AdSn          string `xorm:"not null default '' VARCHAR(255)"`
	AdName        string `xorm:"not null default '' VARCHAR(255)"`
	AdvertiserId  int    `xorm:"not null INT(10)"`
	PayType       int    `xorm:"not null default 0 TINYINT(3)"`
	IntoPrice     int    `xorm:"not null default 0 INT(10)"`
	PutonPlatform int    `xorm:"not null default 1 TINYINT(1)"`
	PutonUrl      string `xorm:"default '' VARCHAR(255)"`
	Remark        string `xorm:"default '' VARCHAR(255)"`
	OwnerApiurl   string `xorm:"default '' VARCHAR(255)"`
	Urlencode     int    `xorm:"not null default 1 TINYINT(2)"`
	StartTime     int    `xorm:"INT(11)"`
	EndTime       int    `xorm:"INT(11)"`
	CreateTime    int    `xorm:"INT(11)"`
	UpdateTime    int    `xorm:"INT(11)"`
}

func GetById(id int) *Advertisement {
	advm := new(Advertisement)
	has, err := engine.Id(id).Get(advm)
	if err != nil || !has {
		log.Println("查找广告失败", id)
		return nil
	}
	return advm
}
