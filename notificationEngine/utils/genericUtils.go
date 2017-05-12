package utils

import (
	"fmt"
	"strconv"
	"time"
)

func GetTimeinSeconds(seconds string) time.Time {
	i, err := strconv.ParseInt(seconds, 10, 64)
	if err != nil {
		panic(err)
	}
	tm := time.Unix(i, 0)
	fmt.Println(tm)
	// const layout = "2017-05-12T09:15:36"
	// fmt.Println(tm.Format(layout))
	// fmt.Println(tm.UTC().Format(layout))
	return tm
}
