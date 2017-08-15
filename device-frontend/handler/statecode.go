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

func illegalData(c *gin.Context) {
	c.JSON(http.StatusBadRequest, gin.H{
		"desc": "The data that you have submitted is not all correct !",
	})
}

func duplicateData(c *gin.Context) {
	c.JSON(http.StatusNotAcceptable, gin.H{
		"desc": "The data that you have submitted seems to be a copy of former data.",
	})
}
