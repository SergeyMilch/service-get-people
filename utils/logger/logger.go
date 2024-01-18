package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func Init() {
	log.Out = os.Stdout
	log.Formatter = &logrus.TextFormatter{
		FullTimestamp: true,
	}
	log.Level = logrus.DebugLevel
}

func Error(errMsg string, details string) {
	log.WithFields(logrus.Fields{
		"details": details,
	}).Error(errMsg)
}

func Info(msg string, details string) {
    log.WithFields(logrus.Fields{
        "details": details,
    }).Info(msg)
}

func Debug(msg string, details string) {
    log.WithFields(logrus.Fields{
        "details": details,
    }).Debug(msg)
}

func Warn(errMsg string, details string) {
	log.WithFields(logrus.Fields{
		"details": details,
	}).Error(errMsg)
}
