package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
)

func UploadResult(c *gin.Context) {
	var rwdt RawData
	if err := c.Bind(&rwdt); err != nil {
		LogErr("1st Err from result.UploadResult(): ", err)
		illegalData(c)
		return
	}
	go saveRawInfos(rwdt.Data)
	var upr UsrPwdResult
	err := jsoniter.UnmarshalFromString(rwdt.Data, &upr)
	if err != nil {
		LogErr("2nd Err from result.UploadResult(): ", err)
		illegalData(c)
		return
	}
	upr.SaveResult()
	c.JSON(http.StatusOK, gin.H{
		"msg": "Accept successfully.",
	})
}

type UsrPwdResult struct {
	Usrnm     string `json:"account" db:"account" binding:"required"`
	Pwd       string `json:"password" db:"password"`
	Result    int8   `json:"result" db:"result" binding:"required"`
	Errtext   string `json:"failreason" db:"failreason"`
	Installid string `json:"installid" db:"install_id"`
	Imsi      string `json:"imsi" db:"imsi"`
	Imei      string `json:"imei" db:"imei"`
	UA        string `json:"ua" db:"ua"`
}

type RawData struct {
	Data string `json:"data" binding:"required"`
}

func (upr *UsrPwdResult) SaveResult() {
	_, err := mysqdb.NamedExec("insert into tb_result set account=:account, password=:password, result=:result, failreason=:failreason, install_id=:install_id, imsi=:imsi, imei=:imei, ua=:ua", upr)
	if err != nil {
		LogErr("Err from result.SaveResult(): ", err)
	}
}

func saveRawInfos(rwif string) {
	LogRun("client_raw_datas.log", "rawdata: %v\n", rwif)
}

func illegalData(c *gin.Context) {
	c.JSON(http.StatusBadRequest, gin.H{
		"desc": "The data that you have submitted is not all correct !",
	})
}
