package tools

import (
	"time"
)

// TimeTrack print duration for specific function
// Example : defer TimeTracker(time.Now(), "CountRecord").
func TimeTrack(startTime time.Time, funcName string) {
	duration := time.Since(startTime)
	InfoLogger.Printf("%s Duration: %s", funcName, duration)
}
