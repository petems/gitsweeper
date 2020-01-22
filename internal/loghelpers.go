package internal

import (
	"os"

	log "github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

func SetupLogger(debug bool) {
	log.SetOutput(os.Stderr)
	textFormatter := new(prefixed.TextFormatter)
	textFormatter.FullTimestamp = true
	textFormatter.TimestampFormat = "01 Jan 2019 15:04:05"
	log.SetFormatter(textFormatter)
	log.SetLevel(log.FatalLevel)

	if debug {
		log.SetLevel(log.InfoLevel)
		log.Info("--debug setting detected - Info level logs enabled")
	}
}

func FatalError(msg string, err error) {
	if err != nil {
		log.WithError(err).Fatal(msg)
	} else {
		log.Fatal(msg)
	}
}
