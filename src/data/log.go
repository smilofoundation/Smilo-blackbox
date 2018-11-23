package data

import (
	"github.com/onrik/logrus/filename"
	"github.com/sirupsen/logrus"
)

var log *logrus.Entry = logrus.WithField("package", "data")

// SetLogger set the logger
func SetLogger(loggers *logrus.Entry) {
	log = loggers.WithFields(log.Data)

	filenameHook := filename.NewHook()

	logrus.AddHook(filenameHook)

}

func init() {
	SetLogger(log)
}
