package utils

import (
	"log"
	"os"
)

var Log *log.Logger

func init() {
	logfile, err := os.OpenFile("log.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
	if err != nil {
		log.Fatal("Failed to get log file:", err)
	}

	Log = log.New(logfile, "", log.Ldate|log.Ltime|log.Lshortfile)
}

func Logger(message string) {
	Log.Println(message)
}
