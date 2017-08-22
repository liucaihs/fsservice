package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func innerErr(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"code": 500,
		"desc": "The server is busy, please try again later.",
	})
}

func illegalData(c *gin.Context) {
	c.JSON(http.StatusBadRequest, gin.H{
		"code": 400,
		"desc": "The data that you have submitted is not all correct !",
	})
}

func rejectData(c *gin.Context) {
	c.JSON(http.StatusNotAcceptable, gin.H{
		"code": 406,
		"desc": "Thanks for your contribution. However, as we have enough data to be used, thus we don not need your data just now. Please devote your data later~",
	})
}

func emptyData(c *gin.Context) {
	c.JSON(http.StatusAlreadyReported, gin.H{
		"code": 208,
		"desc": "Sorry, there is no data that you need currently! Please try again later~",
	})
}
