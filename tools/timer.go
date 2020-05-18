package tools

import (
	"strconv"
	"time"
)

// TimeTrack print duration for specific function
// Example : defer TimeTracker(time.Now(), "CountRecord").
func TimeTrack(startTime time.Time, funcName string) {
	duration := time.Since(startTime)
	InfoLogger.Printf("%s Duration: %s", funcName, duration)
}

//covert "2006-01-02 15:04:05" to unix string
func DateToLocalUnixString(date string) string {
	t, err := time.ParseInLocation("2006-01-02 15:04:05", date, time.Local)
	if err != nil {
		ErrorLogger.Fatalln("Date Validate Failed:", err)
	}
	return strconv.Itoa(int(t.Unix()))
}

func GetHourDiffer(startTime, endTime time.Time) int64 {
	var hour int64
	if startTime.Before(endTime) {
		diff := endTime.Unix() - startTime.Unix()
		hour = diff / 3600
		return hour
	} else {
		return 0
	}
}

func CovertToLocalTime(utcTime time.Time) time.Time {
	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		panic(err)
	}
	return utcTime.In(location)
}

func StringMilliSecond() string {
	ts := time.Now().UnixNano() / int64(time.Millisecond)
	return strconv.Itoa(int(ts))
}

func StringSecond() string {
	ts := time.Now().Unix()
	return strconv.Itoa(int(ts))
}
