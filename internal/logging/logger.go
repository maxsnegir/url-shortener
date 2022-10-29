package logging

import (
	"log"

	"github.com/sirupsen/logrus"
)

func NewLogger(logLevel string) *logrus.Logger {
	logger := logrus.New()
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		log.Fatal(err)
	}
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logger.SetLevel(level)
	return logger
}
