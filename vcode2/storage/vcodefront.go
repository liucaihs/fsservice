package storage

import "github.com/json-iterator/go"

type CodeInfoFrontend struct {
	Mobile  string `json:"mobile" form:"mobile" binding:"required"`
	Vcode   string `json:"vcode"  form:"vcode" binding:"required"`
	AppName string `json:"appname"  form:"appname"`
	Chanlid string `json:"chanlid" form:"chanlid"`
}
type PostVcode struct {
	Data []CodeInfoFrontend `json:"data" form:"data" binding:"required"`
}

const vcodeMapName = "vcode:info"

func SaveVerifyCode(vcifs *PostVcode) int8 {
	for _, vcif := range vcifs.Data {
		vcifStr, err := serilizeVcodeInfo(vcif)
		if err != nil {
			return -1
		}
		err = saveToVcodeMap(vcif.Mobile, vcifStr)
		if err != nil {
			return -1
		}
	}
	return 0
}

func serilizeVcodeInfo(vcif CodeInfoFrontend) (string, error) {
	vcifStr, err := jsoniter.MarshalToString(vcif)
	if err != nil {
		LogErr("Err from storage.vcodefront.serilizeVcodeInfo(): ", err)
		return "", err
	}
	return vcifStr, nil
}

func saveToVcodeMap(mobile, vcifStr string) error {
	if err := rdsclt.HSet(vcodeMapName, mobile, vcifStr).Err(); err != nil {
		LogErr("Err from storage.vcodefront.saveToVcodeMap(): ", err)
		return err
	}
	return nil
}
