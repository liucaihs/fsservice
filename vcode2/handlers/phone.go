package handlers

import (
	"net/http"
	"strconv"
	"vcode2/storage"

	"github.com/gin-gonic/gin"
)

func CollectPhoneNumber(c *gin.Context) {
	var pis storage.PostPhone
	if err := c.Bind(&pis); err != nil {
		storage.LogErr("Err from handlers.phone.CollectPhoneNumber(): ", err)
		illegalData(c)
		return
	}
	runcode, phones := storage.SavePhone(&pis)
	if runcode == -1 {
		innerErr(c)
		return
	}
	if runcode == 2 {
		rejectData(c)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":       200,
		"msg":        "Thanks for your contribution. And as follows, we have used some datas that you provided, please offer their verify codes later.",
		"phonesUsed": phones,
	})
}

func FetchPhoneNumber(c *gin.Context) {
	appnm := c.Param("pkgname")
	cntStr := c.Param("count")
	runcode, pinfs := storage.GetPhone(cntStr, appnm)
	if runcode == 3 {
		illegalData(c)
		return
	}
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
		"data": pinfs.Data,
	})
}

func UpdateSetSize(c *gin.Context) {
	nwszStr := c.Query("newsize")
	nwszInt, err := strconv.Atoi(nwszStr)
	if err != nil {
		storage.LogErr("Err from handlers.phone.UpdateSetSize(): ", err)
		illegalData(c)
		return
	}
	storage.ModifyShudSz(nwszInt)
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": "Update Successfully!",
	})
}
