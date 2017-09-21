package handler

import (
	"net/http"
	"vbs/storage"

	"github.com/gin-gonic/gin"
)

func ObtainOnlinePhoneNumber(c *gin.Context) {
	phoneInfo, rs := storage.GetMobile()
	if rs == "empty" {
		emptyData(c)
		return
	}
	//visit message queue
	c.JSON(http.StatusOK, gin.H{
		"phone": phoneInfo.Mobile,
		"iccid": phoneInfo.Iccid,
		"imsi":  phoneInfo.Imsi,
		"imei":  phoneInfo.Imei,
		"ip":    phoneInfo.IP,
	})
}
