package config

import (
	"io"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	Token       string
	OwnerId     int64
	LoggerId    int64
	DatabaseURI string
	DbName      string
)

var (
	InfoLog  *log.Logger
	ErrorLog *log.Logger
)

func init() {
	_ = godotenv.Load()

	DatabaseURI = os.Getenv("DB_URI")
	DbName = os.Getenv("DB_NAME")
	Token = os.Getenv("TOKEN")
	OwnerId = toInt64(os.Getenv("OWNER_ID"))
	LoggerId = toInt64(os.Getenv("LOGGER_ID"))

	file, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	multi := io.MultiWriter(file, os.Stdout)

	InfoLog = log.New(multi, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLog = log.New(multi, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	setDefaults()
}
