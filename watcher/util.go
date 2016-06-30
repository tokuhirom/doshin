package watcher

import (
	log "github.com/Sirupsen/logrus"
	"time"
)

func ParseDuration(defaultValue time.Duration, userValue string) time.Duration {
	if userValue != "" {
		interval, err := time.ParseDuration(userValue)
		if err != nil {
			log.Fatal("Cannot parse duration: '", userValue, "' ", err)
		}
		return interval
	} else {
		return defaultValue
	}
}
