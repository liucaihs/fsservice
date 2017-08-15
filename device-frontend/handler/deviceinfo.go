package handler

import (
	"device-frontend/storage"

	"github.com/gin-gonic/gin"
)

func DeviceInfoRegister(c *gin.Context) {
	var dvcinf storage.DeviceInfo

	err := c.Bind(&dvcinf)
	if err != nil {
		storage.LogErr("Err from handler.DeviceInfoRegister(): ", err)
		illegalData(c)
		return
	}
	runcode := dvcinf.InsertADeviceInfo()
	if runcode == -1 {
		innerErr(c)
		return
	}
	if runcode == 1 {
		duplicateData(c)
		return
	}
	c.JSON(200, gin.H{
		"msg": "Save successfully~",
	})
}
