package date

import (
	"fmt"
	"strconv"
	"time"
)

const localeStr = "America/Toronto"

// Get today's date formatted in database compatible way
func GetToday() string {
	location, _ := time.LoadLocation(localeStr)

	day := time.Now().In(location).Day()
	month := time.Now().In(location).Month()
	year := time.Now().In(location).Year()
	monthStr := fmt.Sprintf("%02d", int(month))
	dayStr := fmt.Sprintf("%02d", day)
	yearStr := strconv.Itoa(year)
	formattedDate := fmt.Sprintf("%s-%s-%s", yearStr, monthStr, dayStr)

	return formattedDate
}
