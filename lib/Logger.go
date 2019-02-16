package lib

import (
	"fmt"
	"os"
	"log"

)

var logger *log.Logger

func init() {
	file, err := os.OpenFile("logs/logger.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Failed to open log file", err)
	}
	logger = log.New(file, "INFO ", log.Ldate|log.Ltime)
}

// for logging
func Info(args ...interface{}) {
	logger.SetPrefix("[INFO] ")
	logger.Println(args...)
}

func Error(args ...interface{}) {
	logger.SetPrefix("[ERROR] ")
	logger.Println(args...)
}

func Warning(args ...interface{}) {
	logger.SetPrefix("[WARNING] ")
	logger.Println(args...)
}