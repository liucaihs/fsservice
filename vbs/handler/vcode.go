package handler

import (
	"net/http"
	"vbs/storage"

	"github.com/gin-gonic/gin"
)

func SyncProvideIdentifyCode(c *gin.Context) {
	pkg, mobile := c.Param("pkg"), c.Param("phoneNumber")
	vcinf := storage.GetVerifyCode(pkg, mobile)
	if len(vcinf) < 6 {
		emptyData(c)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"phone": mobile,
		"vcode": vcinf,
		"pkg":   pkg,
	})
}
