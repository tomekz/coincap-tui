package ui

import (
	"os"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func init() {
	if os.Getenv("DEBUG") == "true" {
		log.SetLevel(logrus.DebugLevel)
		file, err := os.OpenFile("debug.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
		if err == nil {
			log.Out = file
		} else {
			log.Info("Failed to log to file, using default stderr")
		}
	} else {
		// Only log the warning severity or above.
		log.SetLevel(logrus.PanicLevel)
	}
}
