package storage

import (
	"log"
	"os"
	"strconv"
)

var fdir string

func init() {
	fdir = strconv.Itoa(os.Getpid()) + "/logs/"
}

func mkdirOrdoNothing(filename string) error {
	return os.MkdirAll(fdir, 0755)
}

func LogErr(desc string, err error) {
	if mkdirOrdoNothing(fdir) == nil {
		logfile, _ := os.OpenFile(fdir+"runErr.log", os.O_RDWR|os.O_APPEND|os.O_CREATE, os.ModePerm)
		defer logfile.Close()
		logger := log.New(logfile, "\n", log.Ldate|log.Ltime|log.Lmicroseconds)
		logger.Println(desc, err)
	}
}

func LogRun(filename, format string, v ...interface{}) {
	if mkdirOrdoNothing(fdir) == nil {
		logfile, _ := os.OpenFile(fdir+filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, os.ModePerm)
		defer logfile.Close()
		logger := log.New(logfile, "\n", log.Ldate|log.Ltime|log.Lmicroseconds)
		logger.Printf(format, v)
	}
}
