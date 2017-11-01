package model

type ActivatedIfa struct {
	Id         int64  `xorm:"pk autoincr BIGINT(20)"`
	Ifa        string `xorm:"VARCHAR(255)"`
	Timestamp  int    `xorm:"INT(11)"`
	Status     int    `xorm:"default 0 TINYINT(1)"`
	CreateTime int    `xorm:"INT(11)"`
	FromIp     string `xorm:"not null VARCHAR(255)"`
	RawData    string `xorm:"not null TEXT"`
}

func (this *ActivatedIfa) Insert() (int64, error) {
	return engine.Insert(this)
}
