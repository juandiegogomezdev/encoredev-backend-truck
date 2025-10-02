package utils

import "time"

func GetCurrentTime() time.Time {
	return time.Now().UTC()
}

func GetExpiryTime(duration time.Duration) time.Time {
	return time.Now().UTC().Add(duration)
}
