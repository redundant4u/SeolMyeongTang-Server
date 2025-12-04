package logger

import "log"

func Info(format string) {
	log.Printf("[INFO] " + format)
}

func Error(format string, err error) {
	log.Printf("[ERR] "+format+": %v", err)
}

func Warn(format string) {
	log.Printf("[WARN] " + format)
}

func Fatal(err error, format string) {
	log.Fatalf("[FATAL] "+format+": %v", err)
}
