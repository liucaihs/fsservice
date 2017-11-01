package model

type Advertiser struct {
	Id            int    `xorm:"not null pk autoincr INT(10)"`
	Name          string `xorm:"not null default '' VARCHAR(255)"`
	PayType       int    `xorm:"not null default 0 TINYINT(2)"`
	IntoPrice     string `xorm:"not null default 0.00 DECIMAL(10,2)"`
	PutonPlatform string `xorm:"VARCHAR(255)"`
	Address       string `xorm:"VARCHAR(255)"`
	Requirement   string `xorm:"VARCHAR(500)"`
	Industry      string `xorm:"VARCHAR(255)"`
	Status        int    `xorm:"default 0 TINYINT(1)"`
	CreateTime    int    `xorm:"INT(11)"`
	UpdateTime    int    `xorm:"INT(11)"`
	Contacts      string `xorm:"VARCHAR(255)"`
	Phone         string `xorm:"default '' VARCHAR(20)"`
	Sales         string `xorm:"VARCHAR(255)"`
	Remark        string `xorm:"VARCHAR(255)"`
	RetSucc       string `xorm:"VARCHAR(255)"`
	RetFail       string `xorm:"VARCHAR(255)"`
}
