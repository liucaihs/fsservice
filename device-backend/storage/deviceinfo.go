package storage

import (
	"unsafe"

	"github.com/json-iterator/go"
)

const LISTSIZE = 100

type DeviceInfo struct {
	Id                  int
	Imei                string `json:"imei"`
	Imsi                string `json:"imsi"`
	Phone               string `json:"phone"`
	Simserial           string `json:"simserial"`
	Simcountryiso       string `json:"simcountryiso"`
	Simoperator         string `json:"simoperator"`
	Simoperatorname     string `json:"simoperatorname"`
	Networkcountryiso   string `json:"networkcountryiso"`
	Networkoperator     string `json:"networkoperator"`
	Networkoperatorname string `json:"networkoperatorname"`
	Wifimac             string `json:"wifimac"`
	Bssid               string `json:"bssid"`
	Ssid                string `json:"ssid"`
	Bluemac             string `json:"bluemac"`
	Bluename            string `json:"bluename"`
	Model               string `json:"model"`
	Manufacturer        string `json:"manufacturer"`
	Brand               string `json:"brand"`
	Hardware            string `json:"hardware"`
	Board               string `json:"board"`
	Serial              string `json:"serial"`
	Device              string `json:"device"`
	Buildid             string `json:"buildid"`
	Product             string `json:"product"`
	Display             string `json:"display"`
	Fingerprint         string `json:"fingerprint"`
	Nowrelease          string `json:"nowrelease"`
	Sdk                 string `json:"sdk"`
	Radioversion        string `json:"radioversion"`
	Androidid           string `json:"androidid"`
	Mnc                 string `json:"mnc"`
	Mcc                 string `json:"mcc"`
	Latitude            string `json:"latitude"`
	Longitude           string `json:"longitude"`
	Gsmlac              string `json:"gsmlac"`
	Gsmcid              string `json:"gsmcid"`
	Cdmalatitude        string `json:"cdmalatitude"`
	Cdmalongitude       string `json:"cdmalongitude"`
	Cdmabid             string `json:"cdmabid"`
	Cdmasid             string `json:"cdmasid"`
	Cdmanid             string `json:"cdmanid"`
	Wh                  string `json:"wh"`
	Pkginfos            string `json:"pkginfos"`
}

type AppId struct {
	Applicationid string
	Nextid        int
	Maxreachid    int
	Reset         int8
}

func isEmptyAppList(appname string) (bool, error) {
	rst, err := cache.LLen(appname).Result()
	if err != nil {
		LogErr("Err from storage.deviceinfo.isEmptyAppList(): ", err)
		return true, err
	}
	if rst > 0 {
		return false, nil
	}
	return true, nil
}

func getAInfoFromCache(appname string) *DeviceInfo {
	rstStr, err := cache.LPop(appname).Result()
	if err != nil {
		LogErr("1st Err from storage.deviceinfo.getAInfoFromCache(): ", err)
		return nil
	}
	rstBytes := *(*[]byte)(unsafe.Pointer(&rstStr))
	var dvinf = DeviceInfo{}
	if err = jsoniter.Unmarshal(rstBytes, &dvinf); err != nil {
		LogErr("2nd Err from storage.deviceinfo.getAInfoFromCache(): ", err)
		return nil
	}
	return &dvinf
}

func getInfosFromDB(startId, endId int) []DeviceInfo {
	var dvcinfos = []DeviceInfo{}
	err := db.Select(&dvcinfos, `select * from device_info where id between ? and ? order by id desc`, startId, endId)
	if err != nil {
		LogErr("Err from storage.deviceinfo.getInfosFromDB(): ", err)
		return nil
	}
	return dvcinfos
}

func fillAppList(appname string, dinfos []DeviceInfo) error {
	for _, dinfo := range dinfos {
		dinfstr, err := jsoniter.MarshalToString(dinfo)
		if err != nil {
			LogErr("1st Err from storage.deviceinfo.fillAppList(): ", err)
			return err
		}
		if err := cache.LPush(appname, dinfstr).Err(); err != nil {
			LogErr("2nd Err from storage.deviceinfo.fillAppList(): ", err)
			return err
		}
		//LogRun("tmpdata.log", "dinfstr: %#v\n", dinfstr)
	}
	return nil
}

func getAppidFromDb(appname string) *AppId {
	var ai = AppId{}
	if err := db.Get(&ai, "select * from app_id where applicationid =?", appname); err != nil {
		LogErr("1st Err from storage.deviceinfo.getAppidFromDb(): ", err)
		ai.Applicationid = appname
	}
	var mxid int
	if err := db.Get(&mxid, "select max(id) from device_info"); err != nil {
		LogErr("2nd Err from storage.deviceinfo.getAppidFromDb(): ", err)
		return nil
	}
	if ai.Reset == 1 && ai.Maxreachid <= mxid {
		ai.Nextid = mxid
	}
	return &ai
}

func resetAppQueryId(appname string, ai *AppId) error {
	LogRun("tmpdata.log", "appid(reset): %#v\n", *ai)
	tx, err := db.Beginx()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	if err != nil {
		LogErr("1st Err from storage.deviceinfo.resetAppQueryId(): ", err)
		return err
	}
	query1 := "delete from app_id where applicationid=?"
	_, err = db.Exec(query1, appname)
	if err != nil {
		LogErr("2nd Err from storage.deviceinfo.resetAppQueryId(): ", err)
		return err
	}
	query2 := "insert into app_id values(:applicationid,:nextid,:maxreachid,:reset)"
	_, err = db.NamedExec(query2, ai)
	if err != nil {
		LogErr("3rd Err from storage.deviceinfo.resetAppQueryId(): ", err)
		return err
	}
	if err = tx.Commit(); err != nil {
		LogErr("4th Err from storage.deviceinfo.resetAppQueryId(): ", err)
		return err
	}
	return nil
}

func ReadDeviceInfo(applicName string) (*DeviceInfo, int8) {
	empty, err := isEmptyAppList(applicName)
	if err != nil {
		return nil, -1
	}
	//LogRun("tmpdata.log", "empty: %#v\n", empty)
	if empty {
		appid := getAppidFromDb(applicName)
		if appid == nil {
			return nil, -1
		}
		//LogRun("tmpdata.log", "appid: %#v\n", appid)
		dinfos := getInfosFromDB(appid.Nextid, appid.Nextid+LISTSIZE)
		if dinfos == nil {
			return nil, -1
		}
		//("tmpdata.log", "dinfos: %#v\n", dinfos)
		if len(dinfos) > 0 {
			if err := fillAppList(applicName, dinfos); err != nil {
				return nil, -1
			}
			appid.Nextid = dinfos[0].Id + 1
			appid.Maxreachid = appid.Nextid
			resetAppQueryId(applicName, appid)
		} else {
			appid.Nextid = 1
			appid.Reset = 1
			resetAppQueryId(applicName, appid)
			return nil, 1
		}
	}
	dvinfptr := getAInfoFromCache(applicName)
	if dvinfptr == nil {
		return nil, -1
	}
	//LogRun("tmpdata.log", "*dvinfptr: %#v\n", *dvinfptr)
	return dvinfptr, 0
}
