package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func illegalData(c *gin.Context) {
	c.JSON(http.StatusBadRequest, gin.H{
		"desc": "The data that you have submitted is not all correct !",
	})
}

func emptyData(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{
		"desc": "Sorry, there is no data that you need currently! Please try again later~",
	})
}
