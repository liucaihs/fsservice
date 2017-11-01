package model

type Shortchain struct {
	Id         int    `xorm:"not null pk autoincr INT(11)"`
	Code       string `xorm:"VARCHAR(255)"`
	AdvertUrl  string `xorm:"VARCHAR(255)"`
	Remark     string `xorm:"VARCHAR(255)"`
	CreateTime int    `xorm:"INT(11)"`
	UpdateTime int    `xorm:"INT(11)"`
}

func GetShortBycode(code string) *Shortchain {
	short := new(Shortchain)
	engine.Where("code = ?", code).Get(short)
	return short
}
