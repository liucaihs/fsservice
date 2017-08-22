package storage

import (
	"github.com/go-redis/redis"
	jsoniter "github.com/json-iterator/go"
)

func GetVcode(pkg string) (*PostVcode, int8) {
	var vcifs = PostVcode{}
	var cif CodeInfoFrontend

	<-visitAPP
	phones, ok := appAndItsTels[pkg]
	visitAPP <- 1
	if !ok {
		return nil, 1
	}
	for _, phone := range phones {
		vcifStr, err := rdsclt.HGet(vcodeMapName, phone).Result()
		if err != nil && err != redis.Nil {
			LogErr("1st Err from storage.vcodeback.GetVcode(): ", err)
			return nil, -1
		}
		err = jsoniter.UnmarshalFromString(vcifStr, &cif)
		if err != nil {
			LogErr("2nd Err from storage.vcodeback.GetVcode(): ", err)
			return nil, -1
		}
		vcifs.Data = append(vcifs.Data, cif)
	}
	if err := rdsclt.HDel(vcodeMapName, phones...).Err(); err != nil {
		LogErr("3rd Err from storage.vcodeback.GetVcode(): ", err)
		return nil, -1
	}
	<-visitAPP
	delete(appAndItsTels, pkg)
	visitAPP <- 1
	return &vcifs, 0
}
