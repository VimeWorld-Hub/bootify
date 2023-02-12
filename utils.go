package main

import (
	"fmt"
	"time"
)

func getCurrentDate() int64 {
	return time.Now().Unix()
}

func convertDateMini(date int) string {
	unix := time.Unix(int64(date), 0)
	return fmt.Sprintf("%v.%v.%v", unix.Day(), int(unix.Month()), unix.Year())
}

func convertDateFull(date int) string {
	unix := time.Unix(int64(date), 0)
	return fmt.Sprintf("%v.%v.%v %v:%v:%v", unix.Day(), int(unix.Month()), unix.Year(), unix.Hour(), unix.Minute(), unix.Second())
}

func remove(slice []string, need any) []string {
	for index, value := range slice {
		if value == need {
			return append(slice[:index], slice[index+1:]...)
		}
	}

	return nil
}

func find(array []string, need any) bool {
	for _, value := range array {
		if value == need {
			return true
		}
	}

	return false
}
