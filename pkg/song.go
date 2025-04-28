package pkg

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

var (
	NoAuthTokenError = "Cannot access current song info without YT Desktop api token"
	authHeader       = "authorization"

	ytDesktopApiUrl      = "http://localhost:9863/api/v1/"
	ytDesktopGetStateUrl = ytDesktopApiUrl + "state"

	ytVideoUrl = "https://www.youtube.com/watch?v="
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

func GetCurrentSongInfo(authToken string) string {
	req, err := http.NewRequest("GET", ytDesktopGetStateUrl, nil)
	CheckError(err)

	req.Header.Set(authHeader, authToken)

	client := &http.Client{}
	res, err := client.Do(req)
	CheckError(err)

	resBody, err := io.ReadAll(res.Body)
	CheckError(err)

	var resJson ytState
	json.Unmarshal(resBody, &resJson)

	if resJson.Error != "" {
		return resJson.Error
	}

	return formatCurrentSongInfo(resJson.Video, resJson.Player)
}

func formatCurrentSongInfo(video ytVideo, player ytPlayer) string {
	videoUrl := ytVideoUrl + video.Id
	timestamp := FormatTime(int(player.VideoProgress))

	return fmt.Sprintf("%s %s timestamp: %s", video.Title, videoUrl, timestamp)
}
