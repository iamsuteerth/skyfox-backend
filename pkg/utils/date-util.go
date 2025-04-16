package utils

import "time"

func GetDateFromDateStringDefaultToday(dateStr string) (time.Time, error) {
	if dateStr == "" {
		dateStr = time.Now().Format("2006-01-02")
	}
	return time.Parse("2006-01-02", dateStr)
}
