package pkg

import "fmt"

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
		formattedTime += fmt.Sprintf("%d:", hours)
		seconds = seconds % secondsInHour
	}

	minutes := int(seconds / secondsInMinute)
	if minutes != 0 {
		formattedTime += fmt.Sprintf("%d:", minutes)
		seconds = seconds % secondsInMinute
	}

	if seconds != 0 {
		formattedTime += fmt.Sprintf("%d", seconds)
	}

	return formattedTime
}
