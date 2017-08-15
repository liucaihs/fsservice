package storage

type DeviceInfo struct {
	Imei                string `json:"imei" form:"imei" binding:"required"`
	Imsi                string `json:"imsi" form:"imsi"`
	Phone               string `json:"phone" form:"phone"`
	Simserial           string `json:"simserial" form:"simserial"`
	Simcountryiso       string `json:"simcountryiso" form:"simcountryiso"`
	Simoperator         string `json:"simoperator" form:"simoperator"`
	Simoperatorname     string `json:"simoperatorname" form:"simoperatorname"`
	Networkcountryiso   string `json:"networkcountryiso" form:"networkcountryiso"`
	Networkoperator     string `json:"networkoperator" form:"networkoperator"`
	Networkoperatorname string `json:"networkoperatorname" form:"networkoperatorname"`
	Wifimac             string `json:"wifimac" form:"wifimac"`
	Bssid               string `json:"bssid" form:"bssid"`
	Ssid                string `json:"ssid" form:"ssid"`
	Bluemac             string `json:"bluemac" form:"bluemac"`
	Bluename            string `json:"bluename" form:"bluename"`
	Model               string `json:"model" form:"model"`
	Manufacturer        string `json:"manufacturer" form:"manufacturer"`
	Brand               string `json:"brand" form:"brand"`
	Hardware            string `json:"hardware" form:"hardware"`
	Board               string `json:"board" form:"board"`
	Serial              string `json:"serial" form:"serial"`
	Device              string `json:"device" form:"device"`
	Buildid             string `json:"buildid" form:"buildid"`
	Product             string `json:"product" form:"product"`
	Display             string `json:"display" form:"display"`
	Fingerprint         string `json:"fingerprint" form:"fingerprint"`
	Nowrelease          string `json:"nowrelease" form:"nowrelease"`
	Sdk                 string `json:"sdk" form:"sdk"`
	Radioversion        string `json:"radioversion" form:"radioversion"`
	Androidid           string `json:"androidid" form:"androidid"`
	Mnc                 string `json:"mnc" form:"mnc"`
	Mcc                 string `json:"mcc" form:"mcc"`
	Latitude            string `json:"latitude" form:"latitude"`
	Longitude           string `json:"longitude" form:"longitude"`
	Gsmlac              string `json:"gsmlac" form:"gsmlac"`
	Gsmcid              string `json:"gsmcid" form:"gsmcid"`
	Cdmalatitude        string `json:"cdmalatitude" form:"cdmalatitude"`
	Cdmalongitude       string `json:"cdmalongitude" form:"cdmalongitude"`
	Cdmabid             string `json:"cdmabid" form:"cdmabid"`
	Cdmasid             string `json:"cdmasid" form:"cdmasid"`
	Cdmanid             string `json:"cdmanid" form:"cdmanid"`
	Wh                  string `json:"wh" form:"wh"`
	Pkginfos            string `json:"pkginfos" form:"pkginfos"`
}

var receivedDid = make(map[string]int8)
var allowedVisit = make(chan int8, 1)

func (dvif *DeviceInfo) InsertADeviceInfo() int8 {
	<-allowedVisit
	if receivedDid[dvif.Imei] == 1 {
		allowedVisit <- 1
		return 1
	}
	allowedVisit <- 1

	tx, err := db.Beginx()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	if err != nil {
		LogErr("1st Err from storage.deviceinfo.InsertADeviceInfo(): ", err)
		return -1
	}
	query := `INSERT INTO device_info(imei, imsi, phone, simserial, simcountryiso, simoperator, simoperatorname, networkcountryiso, networkoperator, networkoperatorname, wifimac, bssid, ssid, bluemac, bluename, model, manufacturer, brand, hardware, board, serial, device, buildid, product, display, fingerprint, nowrelease, sdk, radioversion, androidid, mnc, mcc, latitude, longitude, gsmlac, gsmcid, cdmalatitude, cdmalongitude, cdmabid, cdmasid, cdmanid, wh, pkginfos) VALUES (:imei, :imsi, :phone, :simserial, :simcountryiso, :simoperator, :simoperatorname, :networkcountryiso, :networkoperator, :networkoperatorname, :wifimac, :bssid, :ssid, :bluemac, :bluename, :model, :manufacturer, :brand, :hardware, :board, :serial, :device, :buildid, :product, :display, :fingerprint, :nowrelease, :sdk, :radioversion, :androidid, :mnc, :mcc, :latitude, :longitude, :gsmlac, :gsmcid, :cdmalatitude, :cdmalongitude, :cdmabid, :cdmasid, :cdmanid, :wh, :pkginfos)`
	_, err = tx.NamedExec(query, dvif)
	if err != nil {
		LogErr("2nd Err from storage.deviceinfo.InsertADeviceInfo(): ", err)
		return -1
	}
	if err = tx.Commit(); err != nil {
		LogErr("3rd Err from storage.deviceinfo.InsertADeviceInfo(): ", err)
		return -1
	}

	<-allowedVisit
	receivedDid[dvif.Imei] = 1
	allowedVisit <- 1
	return 0
}

func imeisInit() error {
	allowedVisit <- 1

	var dids []string
	err := db.Select(&dids, "select imei from device_info")
	if err != nil {
		LogErr("Err from storage.deviceinfo.imeisInit(): ", err)
		return err
	}
	<-allowedVisit
	for _, did := range dids {
		receivedDid[did] = 1
	}
	allowedVisit <- 1
	return nil
}
