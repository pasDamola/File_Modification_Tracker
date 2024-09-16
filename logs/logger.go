package logs

import (
	"log"
	"os"
)

var logger *log.Logger

func InitLogger() {
	file, err := os.OpenFile("tracker.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}

	logger = log.New(file, "", log.Ldate|log.Ltime|log.Lshortfile)
}

func Log(message string) {
	logger.Println(message)
}
