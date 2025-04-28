package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
)

var (
	authHeader string = "authorization"

	ytDesktopApiUrl      string = "http://localhost:9863/api/v1/"
	ytDesktopGetStateUrl string = ytDesktopApiUrl + "state"

	ytVideoUrl string = "https://www.youtube.com/watch?v="

	secondsInHour   = 3600
	secondsInMinute = 60
)

type ytState struct {
	Video  ytVideo
	Player ytPlayer

	Error string
}

type ytVideo struct {
	Title string
	Id    string
}

type ytPlayer struct {
	VideoProgress float32
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func formatTime(seconds int) string {
	formattedTime := ""

	hours := int(seconds / secondsInHour)
	if hours != 0 {
		formattedTime += fmt.Sprintf("%dч ", hours)
		seconds = seconds % secondsInHour
	}

	minutes := int(seconds / secondsInMinute)
	if minutes != 0 {
		formattedTime += fmt.Sprintf("%dм ", minutes)
		seconds = seconds % secondsInMinute
	}

	if seconds != 0 {
		formattedTime += fmt.Sprintf("%dс", seconds)
	}

	return formattedTime
}

func getCurrentSongInfo(video ytVideo, player ytPlayer) string {
	videoUrl := ytVideoUrl + video.Id
	timestamp := formatTime(int(player.VideoProgress))

	return fmt.Sprintf("%s %s таймстемп: %s", video.Title, videoUrl, timestamp)
}

func main() {
	authToken := flag.String("token", "", "YT Desktop API token")
	flag.Parse()

	if *authToken == "" {
		fmt.Println("Не могу получить информацию о текущей песне без токена авторизации")
		return
	}

	req, err := http.NewRequest("GET", ytDesktopGetStateUrl, nil)
	checkError(err)

	req.Header.Set(authHeader, *authToken)

	client := &http.Client{}
	res, err := client.Do(req)
	checkError(err)

	resBody, err := io.ReadAll(res.Body)
	checkError(err)

	var resJson ytState
	json.Unmarshal(resBody, &resJson)

	if resJson.Error != "" {
		fmt.Println(resJson.Error)
		return
	}

	fmt.Println(getCurrentSongInfo(resJson.Video, resJson.Player))
}
