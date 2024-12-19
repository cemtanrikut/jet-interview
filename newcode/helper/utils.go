package helper

import (
	"io/ioutil"
	"log"
	"time"
)

const lastRunFile = "last_run.txt"

// GetLastRunTime retrieves the last run time from a file.
func GetLastRunTime() time.Time {
	data, err := ioutil.ReadFile(lastRunFile)
	if err != nil {
		return time.Now().Add(-24 * time.Hour) // Default to 24 hours ago
	}

	lastRunTime, err := time.Parse(time.RFC3339, string(data))
	if err != nil {
		return time.Now().Add(-24 * time.Hour)
	}

	return lastRunTime
}

// UpdateLastRunTime updates the last run time in a file.
func UpdateLastRunTime(currentTime time.Time) {
	err := ioutil.WriteFile(lastRunFile, []byte(currentTime.Format(time.RFC3339)), 0644)
	if err != nil {
		log.Printf("Son çalıştırma zamanı güncellenemedi: %v", err)
	}
}
