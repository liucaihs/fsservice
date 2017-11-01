package model

type IfaLog struct {
	Id         int64  `xorm:"pk autoincr BIGINT(20)"`
	Mac        string `xorm:"not null VARCHAR(255)"`
	Ifa        string `xorm:"VARCHAR(255)"`
	Ip         string `xorm:"VARCHAR(255)"`
	RawData    string `xorm:"not null TEXT"`
	CreateTime int    `xorm:"INT(11)"`
}
