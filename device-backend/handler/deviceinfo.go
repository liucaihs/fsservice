package handler

import (
	"device-backend/storage"

	"github.com/gin-gonic/gin"
)

func DeviceInfoDisplay(c *gin.Context) {
	did := c.Param("applicationame")

	dvcinf, runcode := storage.ReadDeviceInfo(did)
	if runcode == -1 {
		innerErr(c)
		return
	}
	if runcode == 1 {
		emptyData(c)
		return
	}
	c.JSON(200, gin.H{
		"data": *dvcinf,
	})
}
