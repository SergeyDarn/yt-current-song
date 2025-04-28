package pkg

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/charmbracelet/lipgloss"
)

var (
	NoAuthTokenError = "Cannot access current song info without YT Desktop api token"
	authHeader       = "authorization"

	ytDesktopApiUrl      = "http://localhost:9863/api/v1/"
	ytDesktopGetStateUrl = ytDesktopApiUrl + "state"

	ytVideoUrl                = "https://www.youtube.com/watch?v="
	ytVideoTimeQuery          = "&t="
	ytStatePaused             = 0
	ytStatePlaying            = 1
	songCollectionMinuteStart = 15
)

var (
	PAUSE_COLOR  = lipgloss.Color("192")
	SONG_COLOR   = lipgloss.Color("#f98b6c")
	ARTIST_COLOR = lipgloss.Color("#ef4fa6")
	URL_COLOR    = lipgloss.Color("#1e88e5")
)

type ytState struct {
	Video  ytVideo
	Player ytPlayer

	Error string
}

type ytVideo struct {
	Author          string
	Title           string
	Id              string
	DurationSeconds int
}

type ytPlayer struct {
	VideoProgress float32
	TrackState    int
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
	if player.TrackState == ytStatePaused {
		return PrepareColorOutput("No song is currently playing", PAUSE_COLOR)
	}

	videoUrl := ytVideoUrl + video.Id
	if isSongCollection(video.DurationSeconds) {
		videoUrl += ytVideoTimeQuery + strconv.Itoa(int(player.VideoProgress))
	}

	return fmt.Sprintf(
		"%s %s %s",
		PrepareColorOutput(video.Author, ARTIST_COLOR),
		PrepareColorOutput(video.Title, SONG_COLOR),
		PrepareColorOutput(videoUrl, URL_COLOR),
	)
}

func isSongCollection(durationSeconds int) bool {
	durationMinutes := SecondsToMinutes(durationSeconds)

	return durationMinutes >= songCollectionMinuteStart
}
