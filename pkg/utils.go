package pkg

import (
	"bytes"
	"io"
	"net/http"
	"strconv"
)

const (
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

func doRequest(method string, url string, body []byte, headers map[string]string) []byte {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	CheckError(err)

	if method == http.MethodPost {
		req.Header.Set("Content-Type", "application/json")
	}

	for name, value := range headers {
		req.Header.Set(name, value)
	}

	client := &http.Client{}
	res, err := client.Do(req)
	CheckError(err)

	resBody, err := io.ReadAll(res.Body)
	CheckError(err)

	return resBody
}
