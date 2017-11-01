package model

import (
	"log"
)

type Setting struct {
	Id         int    `xorm:"not null pk autoincr INT(10)"`
	Key        string `xorm:"not null VARCHAR(255)"`
	Value      string `xorm:"default '' VARCHAR(255)"`
	CreateTime int    `xorm:"INT(11)"`
	UpdateTime int    `xorm:"INT(11)"`
}

func GetSetVal(key string) string {
	setRec := new(Setting)
	has, err := engine.Where("`key`=?", key).Get(setRec)
	if err != nil || !has || len(setRec.Value) == 0 {
		log.Println("查找配置失败:", key, setRec)
		return ""
	}
	return setRec.Value
}
