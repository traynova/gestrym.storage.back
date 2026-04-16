package utils

import "time"

/*
GetCurrentTime
Get the current time
Return:

	time.Time

Example:

	2021-08-25 15:00:00 +0000 UTC
*/
func GetCurrentTime() time.Time {
	return time.Now()
}

/*
GetCurrentTimeFormated
Get the current time formated
Return:

	string

Example:

	"2021-08-25"
*/
func GetCurrentTimeFormated() string {
	currentTime := time.Now()
	return currentTime.Format("2006-01-02")
}
