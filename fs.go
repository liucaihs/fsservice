package main

import (
	"net/http"
	"gopkg.in/gin-gonic/gin.v1"
	"gopkg.in/mgo.v2"
	"io/ioutil"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"strings"
)

const url = "127.0.0.1"
//const url = "192.168.1.211"

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func calcmd5(data []byte) (string) {
	hash := md5.New()
	hash.Write(data)
	cipherText := hash.Sum(nil)
	return hex.EncodeToString(cipherText)
}

func upload(c *gin.Context) {
	name := c.Param("name")
	md5 := c.Param("md5")
	
	data, errR := ioutil.ReadAll(c.Request.Body)
	check(errR)
	defer c.Request.Body.Close()
	rmd5 := calcmd5(data)
	if strings.Compare(md5, rmd5) != 0 {
		c.String(http.StatusUnauthorized, "unmatched md5")
		return
	}

	session, errD := mgo.Dial(url)
	check(errD)
	defer session.Close()

	db := session.DB("storage")
	file, err := db.GridFS("fs").Open(name)
	if err != nil {
		fileNew, errN := db.GridFS("fs").Create(name)
		check(errN)
	
		_, err = fileNew.Write(data)
		check(err)

		fileNew.Close()

		c.Status(http.StatusOK)
	} else {
		defer file.Close()
		c.Status(http.StatusFound)
	}
}

func download(c *gin.Context) {
	name := c.Param("name")

	session, err := mgo.Dial(url)
	check(err)
	defer session.Close()

	db := session.DB("storage")
	file, err := db.GridFS("fs").Open(name)
	if err != nil {
		c.Status(http.StatusNotFound)
	} else {
		var bb bytes.Buffer
		bb.ReadFrom(file)
		defer file.Close()
		c.Data(http.StatusOK, "application/octet-stream", bb.Bytes())
	}
}

func isExist(c *gin.Context) {
	name := c.Param("name")
	session, err := mgo.Dial(url)
	check(err)
	defer session.Close()

	db := session.DB("storage")
	file, err := db.GridFS("fs").Open(name)

	if err != nil {
		c.Status(http.StatusNotFound)
	} else {
		defer file.Close()
		c.JSON(http.StatusOK, gin.H{
			"length": file.Size(),
			"md5": file.MD5(),
		})
	}
}

func main() {
	//gin.SetMode(gin.ReleaseMode)
	
	router := gin.New()

	router.Use(gin.Logger())

	v := router.Group("/v2")
	{
		v.POST("/upload/:name/:md5", upload)
		v.GET("/download/:name", download)
		v.GET("/isExist/:name", isExist)
	}

	router.Run(":8400")
}
