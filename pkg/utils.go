package pkg

import (
	"strconv"
)

var (
	secondsInHour   = 3600
	secondsInMinute = 60
)

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

func FormatTime(seconds int) string {
	formattedTime := ""

	hours := int(seconds / secondsInHour)
	if hours != 0 {
		formattedTime += addLeadingZero(hours) + ":"
		seconds = seconds % secondsInHour
	}

	minutes := int(seconds / secondsInMinute)
	formattedTime += addLeadingZero(minutes) + ":"
	seconds = seconds % secondsInMinute

	formattedTime += addLeadingZero(seconds)

	return formattedTime
}

func addLeadingZero(number int) string {
	strNumber := strconv.Itoa(number)

	if number > 10 {
		return strNumber
	}

	return "0" + strNumber
}
