package model

type ShortLog struct {
	Id         int64  `xorm:"pk autoincr BIGINT(20)"`
	ShortId    int    `xorm:"not null INT(32)"`
	Ip         string `xorm:"VARCHAR(255)"`
	Status     int    `xorm:"default 1 TINYINT(2)"`
	RawData    string `xorm:"not null TEXT"`
	CreateTime int    `xorm:"INT(11)"`
}

func (this *ShortLog) Insert() (int64, error) {
	return engine.Insert(this)
}
