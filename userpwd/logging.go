package main

import (
	"log"
	"os"
	"time"
)

var fdir string

func mkdirOrdoNothing() error {
	fdir = time.Now().Format("2006-01-02") + "/"
	return os.MkdirAll(fdir, 0777)
}

func LogErr(desc string, err error) {
	if mkdirOrdoNothing() == nil {
		logfile, _ := os.OpenFile(fdir+"runErr.log", os.O_RDWR|os.O_APPEND|os.O_CREATE, os.ModePerm)
		defer logfile.Close()
		logger := log.New(logfile, "\n", log.Ldate|log.Ltime|log.Lmicroseconds)
		logger.Println(desc, err)
	}
}

func LogRun(filename, format string, v ...interface{}) {
	if mkdirOrdoNothing() == nil {
		logfile, _ := os.OpenFile(fdir+filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, os.ModePerm)
		defer logfile.Close()
		logger := log.New(logfile, "\n", log.Ldate|log.Ltime|log.Lmicroseconds)
		logger.Printf(format, v)
	}
}
