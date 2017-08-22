package handlers

import (
	"net/http"
	"vcode2/storage"

	"github.com/gin-gonic/gin"
)

func CollectVerifyCode(c *gin.Context) {
	var vcifs storage.PostVcode
	if err := c.Bind(&vcifs); err != nil {
		storage.LogErr("Err from handlers.vcode.CollectVerifyCode(): ", err)
		illegalData(c)
		return
	}
	runcode := storage.SaveVerifyCode(&vcifs)
	if runcode == -1 {
		innerErr(c)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "Save successfully~",
	})
}

func FetchVerifyCode(c *gin.Context) {
	appname := c.Param("pkgname")
	vcifs, runcode := storage.GetVcode(appname)
	if runcode == -1 {
		innerErr(c)
		return
	}
	if runcode == 1 {
		emptyData(c)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": vcifs.Data,
	})
}
