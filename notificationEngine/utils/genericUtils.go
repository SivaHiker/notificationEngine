package utils

import (
	"fmt"
	"time"
)

func GetTimeinSeconds(seconds string) time.Time {
	// i, err := strconv.ParseInt(seconds, 10, 64)
	// if err != nil {
	// panic(err)
	// }
	// tm := time.Unix(i, 0)
	// fmt.Println(tm)
	layout := "2006-01-02T15:04:05"
	t, err := time.Parse(layout, seconds)

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(t)
	return t
}
