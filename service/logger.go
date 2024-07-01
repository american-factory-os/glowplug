package service

import (
	"log"
	"os"
)

func NewLogger() *log.Logger {
	return log.New(os.Stdout, "", log.LstdFlags)
}
