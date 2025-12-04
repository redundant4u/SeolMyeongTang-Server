package logger

import "log"

func Info(format string, args ...any) {
	log.Printf("[INFO] "+format, args...)
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
