package main

import "time"

func getCurrentDate() int64 {
	return time.Now().Unix()
}
