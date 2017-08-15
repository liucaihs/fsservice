package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func innerErr(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"desc": "The server is busy, please try again later.",
	})
}

func emptyData(c *gin.Context) {
	c.JSON(http.StatusAlreadyReported, gin.H{
		"desc": "Nowdays, you have obtained whole data that we have collected. If you want to get them again, please try again later~",
	})
}
