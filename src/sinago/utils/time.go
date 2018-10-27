package utils

import "time"

func GetDate() string {
	return time.Now().Format("2006-01-02")
}
