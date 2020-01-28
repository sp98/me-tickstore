package utility

import (
	"fmt"
	"strings"
	"time"
)

//GetFilteredStocks returns formated 2-D array of stocks
func GetFilteredStocks(stock string, tokens []string) [][]string {
	result := [][]string{}
	stocks := strings.Split(stock, ",")

	for _, s := range stocks {
		sSlice := strings.Split(s, ";")
		if contains(tokens, sSlice[2]) {
			result = append(result, sSlice)
		}
	}

	return result
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

//IsWithInTimeRange checks if the current time period is within a time range.
func IsWithInTimeRange(time1, time2 string) (bool, error) {
	loc, _ := time.LoadLocation("Asia/Kolkata")
	t1String := fmt.Sprintf(time1, time.Now().Format("2006-01-02"))
	t1, err := time.ParseInLocation("2006-01-02 15:04:05", t1String, loc)
	if err != nil {
		return false, fmt.Errorf("error parsing market open time. %+v", err)
	}

	t2String := fmt.Sprintf(time2, time.Now().Format("2006-01-02"))
	t2, err := time.ParseInLocation("2006-01-02 15:04:05", t2String, loc)
	if err != nil {
		return false, fmt.Errorf("error parsing market open time. %+v", err)
	}

	currentTime := time.Now()
	if currentTime.After(t1) && currentTime.Before(t2) && int(currentTime.Weekday()) != 6 && int(currentTime.Weekday()) != 0 {
		return true, nil
	}
	return false, nil

}
