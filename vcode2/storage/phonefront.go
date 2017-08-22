package storage

import (
	"sync"
	"time"

	"github.com/go-redis/redis"
	"github.com/json-iterator/go"
)

type PhoneInfoFrontend struct {
	Mobile string `json:"mobile" form:"mobile" binding:"required"`
	Cid    string `json:"cid" form:"cid"`
	Imsi   string `json:"imsi" form:"imsi"`
	Imei   string `json:"imei" form:"imei"`
	Iccid  string `json:"iccid" form:"iccid"`
	Ip     string `json:"ip" form:"ip"`
}
type PostPhone struct {
	Data []PhoneInfoFrontend `json:"data" form:"data" binding:"required"`
}

const phoneSetName = "mobile:number"
const phoneMapName = "mobile:info"

var shouldSize int
var visitSS = make(chan int8, 1)
var visitRDS = make(chan int8, 1)

func init() {
	shouldSize = 1000
	visitSS <- 1
	visitRDS <- 1
}

func SavePhone(pis *PostPhone) (int8, []string) {
	var mobilesNeeded = []PhoneInfoFrontend{}
	var wg sync.WaitGroup
	var phonesResponsed []string

	<-visitRDS
	curtSize, err := getCurrentCacheSize()
	if err != nil {
		visitRDS <- 1
		return -1, nil
	}
	<-visitSS
	more := shouldSize - curtSize
	visitSS <- 1
	if more <= 0 {
		visitRDS <- 1
		return 2, nil
	}
	if more >= len(pis.Data) {
		mobilesNeeded = pis.Data
	} else {
		mobilesNeeded = pis.Data[:more]
	}

	wg.Add(3)
	go func() {
		defer wg.Done()
		for _, pif := range mobilesNeeded {
			pifStr, _ := serilizePhoneInfo(pif)
			saveToPhoneMap(pif.Mobile, pifStr)
		}
	}()
	go func() {
		defer wg.Done()
		for _, pif := range mobilesNeeded {
			saveToPhonSortedeSet(pif.Mobile)
		}
		visitRDS <- 1
	}()
	go func() {
		defer wg.Done()
		for _, pif := range mobilesNeeded {
			phonesResponsed = append(phonesResponsed, pif.Mobile)
		}
	}()
	wg.Wait()
	return 0, phonesResponsed
}

func getCurrentCacheSize() (int, error) {
	rst, err := rdsclt.ZCard(phoneSetName).Result()
	if err != nil {
		LogErr("Err from storage.phonefront.getCurrentCacheSize(): ", err)
		return -1, err
	}
	return int(rst), nil
}

func serilizePhoneInfo(pif PhoneInfoFrontend) (string, error) {
	pifstr, err := jsoniter.MarshalToString(pif)
	if err != nil {
		LogErr("Err from storage.phonefront.serilizePhoneInfo(): ", err)
		return "", err
	}
	return pifstr, nil
}

func saveToPhoneMap(mobile, pif string) error {
	if err := rdsclt.HSet(phoneMapName, mobile, pif).Err(); err != nil {
		LogErr("Err from storage.phonefront.saveToPhoneMap(): ", err)
		return err
	}
	return nil
}

func saveToPhonSortedeSet(mobile string) error {
	date := time.Now().UnixNano()
	z := redis.Z{Score: float64(date), Member: mobile}
	if err := rdsclt.ZAdd(phoneSetName, z).Err(); err != nil {
		LogErr("Err from storage.phonefront.saveToPhonSortedeSet(): ", err)
		return err
	}
	return nil
}

func ModifyShudSz(nwsz int) {
	<-visitSS
	shouldSize = nwsz
	visitSS <- 1
}
