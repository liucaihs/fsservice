package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SendAccount(c *gin.Context) {
	actinfs, ok := bscElems.Load(mblcach)
	if !ok || len(actinfs.([]Account)) < 3 {
		acts := getAvailableAccounts()
		if acts != nil {
			bscElems.Store(mblcach, acts)
		} else {
			emptyData(c)
			return
		}
	}
	actinfs, _ = bscElems.Load(mblcach)
	usr := actinfs.([]Account)[0]
	bscElems.Store(mblcach, actinfs.([]Account)[1:])
	c.JSON(http.StatusOK, gin.H{
		"username": usr.Usrnm,
		"password": usr.Pwd,
	})
}

type Account struct {
	Usrnm string `db:"account"`
	Pwd   string `db:"password"`
}

func getAvailableAccounts() []Account {
	var acts = []Account{}
	query := "select account, password from view_unused limit 300"
	if err := mysqdb.Select(&acts, query); err != nil {
		LogErr("1st Err from account.getAvailableAccounts(): ", err)
	}
	if len(acts) < 1 {
		return nil
	}
	return acts
}

func emptyData(c *gin.Context) {
	c.JSON(http.StatusGone, gin.H{
		"desc": "Sorry, there is no data that you need currently! Please try again later~",
	})
}
