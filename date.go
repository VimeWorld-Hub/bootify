package main

import (
	"fmt"
	"time"
)

func getCurrentDate() int64 {
	return time.Now().Unix()
}

func convertDate(date int) string {
	unix := time.Unix(int64(date), 0)
	return fmt.Sprintf("%v.%v.%v", unix.Day(), int(unix.Month()), unix.Year())
}

func getCurrentDateFormatted() string {
	now := time.Now()
	return fmt.Sprintf("%v.%v.%v", now.Day(), int(now.Month()), now.Year())
}
