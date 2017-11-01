package model

type Media struct {
	Id          int    `xorm:"not null pk autoincr INT(10)"`
	Name        string `xorm:"not null VARCHAR(255)"`
	Belongs     string `xorm:"VARCHAR(255)"`
	ChannelType string `xorm:"not null default '0' VARCHAR(255)"`
	Contacts    string `xorm:"VARCHAR(255)"`
	Phone       string `xorm:"default '' VARCHAR(20)"`
	Status      int    `xorm:"default 0 TINYINT(1)"`
	CreateTime  int    `xorm:"INT(11)"`
	UpdateTime  int    `xorm:"INT(11)"`
}
