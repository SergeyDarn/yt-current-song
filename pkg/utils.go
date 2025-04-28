package pkg

import (
	"strconv"

	"github.com/charmbracelet/lipgloss"
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

	hours := SecondsToHours(seconds)
	if hours != 0 {
		formattedTime += addLeadingZero(hours) + ":"
		seconds = seconds % secondsInHour
	}

	minutes := SecondsToMinutes(seconds)
	formattedTime += addLeadingZero(minutes) + ":"
	seconds = seconds % secondsInMinute

	formattedTime += addLeadingZero(seconds)

	return formattedTime
}

func SecondsToHours(seconds int) int {
	return int(seconds / secondsInHour)
}

func SecondsToMinutes(seconds int) int {
	return int(seconds / secondsInMinute)
}

func addLeadingZero(number int) string {
	strNumber := strconv.Itoa(number)

	if number >= 10 {
		return strNumber
	}

	return "0" + strNumber
}

func PrepareColorOutput(output string, color lipgloss.Color) string {
	return lipgloss.NewStyle().Foreground(lipgloss.TerminalColor(color)).Render(output)
}
