package storage

import (
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"github.com/json-iterator/go"
)

var appAndItsTels = make(map[string][]string)
var visitAPP = make(chan int8, 1)

func init() {
	visitAPP <- 1
}

func GetPhone(num, pkg string) (int8, *PostPhone) {
	var phones = []string{}
	var phoneInfos = PostPhone{}

	numInt, err := strconv.Atoi(num)
	if err != nil {
		LogErr("Err from storage.phoneback.GetPhone(): ", err)
		return 3, nil
	}
	phones, err = getNewerPhones(numInt)
	if err != nil {
		return -1, nil
	}
	if len(phones) < 1 {
		return 1, nil
	}
	<-visitRDS
	remvPhonesFromSet(phones)
	visitRDS <- 1
	phoneInfos, err = getCorrespondPhoneInfo(phones)
	if err != nil {
		return -1, nil
	}
	err = remvPhonesFromMap(phones)
	if err != nil {
		return -1, nil
	}
	<-visitAPP
	appAndItsTels[pkg] = phones
	visitAPP <- 1
	return 0, &phoneInfos
}

func getNewerPhones(limit int) ([]string, error) {
	var phones = []string{}
	var err error

	opt := redis.ZRangeBy{Min: "(0", Max: strconv.Itoa(int(time.Now().UnixNano())), Count: int64(limit)}
	phones, err = rdsclt.ZRevRangeByScore(phoneSetName, opt).Result()
	if err != nil {
		LogErr("Err from storage.phoneback.getNewerPhones(): ", err)
		return nil, err
	}
	return phones, nil
}

func getCorrespondPhoneInfo(phones []string) (PostPhone, error) {
	var pinfs = PostPhone{}

	for _, mobile := range phones {
		pinfStr, err := rdsclt.HGet(phoneMapName, mobile).Result()
		if err != nil && err != redis.Nil {
			LogErr("Err from storage.phoneback.getCorrespondPhoneInfo(): ", err)
			return pinfs, err
		}
		pinf, err := deserializePhoneInfo(pinfStr)
		if err != nil {
			return pinfs, err
		}
		pinfs.Data = append(pinfs.Data, pinf)
	}
	return pinfs, nil
}

func deserializePhoneInfo(pinf string) (PhoneInfoFrontend, error) {
	var pif = PhoneInfoFrontend{}
	err := jsoniter.UnmarshalFromString(pinf, &pif)
	if err != nil {
		LogErr("Err from storage.phoneback.deserializePhoneInfo(): ", err)
		return pif, err
	}
	return pif, nil
}

func remvPhonesFromSet(phones []string) error {
	for _, phone := range phones {
		if err := rdsclt.ZRem(phoneSetName, phone).Err(); err != nil {
			LogErr("Err from storage.phoneback.remvPhonesFromSet(): ", err)
			return err
		}
	}
	return nil
}

func remvPhonesFromMap(phones []string) error {
	if err := rdsclt.HDel(phoneMapName, phones...).Err(); err != nil {
		LogErr("Err from storage.phoneback.remvPhonesFromMap(): ", err)
		return err
	}
	return nil
}
