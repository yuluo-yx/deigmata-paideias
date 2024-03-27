package log

import (
	"log"
)

type XLogger struct {
	Debug bool
}

// Debugf Log to stdout only if Debug is true.
func (logger XLogger) Debugf(format string, args ...interface{}) {
	if logger.Debug {
		log.Printf(format+"\n", args...)
	}
}

// Infof Log to stdout only if Debug is true.
func (logger XLogger) Infof(format string, args ...interface{}) {
	if logger.Debug {
		log.Printf(format+"\n", args...)
	}
}

// Warnf Log to stdout always.
func (logger XLogger) Warnf(format string, args ...interface{}) {
	log.Printf(format+"\n", args...)
}

// Errorf Log to stdout always.
func (logger XLogger) Errorf(format string, args ...interface{}) {
	log.Printf(format+"\n", args...)
}
