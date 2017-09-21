package storage

import (
	"github.com/json-iterator/go"
)

type PhoneInfo struct {
	Mobile string `json:"mobile"`
	Iccid  string `json:"iccid"`
	Imsi   string `json:"imsi"`
	Imei   string `json:"imei"`
	IP     string `json:"ip"`
}

const MOBILEMAP string = "phone:info"
const MOBILESET string = "phones"

func GetMobile() (*PhoneInfo, string) {
	mobile := getAPhoneFromSet()
	if mobile == "" {
		return nil, "empty"
	}
	mobileInfoptr := getPhoneInfoFromMap(mobile)
	if mobileInfoptr == nil {
		return nil, "empty"
	}
	return mobileInfoptr, "ok"
}

func getAPhoneFromSet() string {
	phone, err := rdsclt.SRandMember(MOBILESET).Result()
	if err != nil {
		LogErr("Err from storage.phone.getAPhoneFromSet(): ", err)
		return ""
	}
	return phone
}

func getPhoneInfoFromMap(phone string) *PhoneInfo {
	pinfStr, err := rdsclt.HGet(MOBILEMAP, phone).Result()
	if err != nil {
		LogErr("1st Err from storage.phone.getPhoneInfoFromMap(): ", err)
		return nil
	}
	var pif = PhoneInfo{}
	if err := jsoniter.UnmarshalFromString(pinfStr, &pif); err != nil {
		LogErr("2nd Err from storage.phone.getPhoneInfoFromMap(): ", err)
		return nil
	}
	return &pif
}
