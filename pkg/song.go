package pkg

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

const (
	NoAuthTokenError = "Cannot access current song info without YT Desktop api token"
	authHeader       = "authorization"

	ytDesktopApiUrl      = "http://localhost:9863/api/v1/"
	ytDesktopGetStateUrl = ytDesktopApiUrl + "state"

	ytVideoUrl                = "https://www.youtube.com/watch?v="
	ytTimeParam               = "&t="
	ytPlaylistParam           = "&list="
	ytStatePaused             = 0
	ytStatePlaying            = 1
	songCollectionMinuteStart = 15
)

type ytVideoState struct {
	Video      ytVideo
	Player     ytPlayer
	PlaylistId string

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
	songState, err := getYtVideoState(authToken)
	if err != nil {
		return err.Error()
	}

	return formatCurrentSongInfo(songState)
}

func getYtVideoState(authToken string) (ytVideoState, error) {
	req, err := http.NewRequest("GET", ytDesktopGetStateUrl, nil)
	CheckError(err)

	req.Header.Set(authHeader, authToken)

	client := &http.Client{}
	res, err := client.Do(req)
	CheckError(err)

	resBody, err := io.ReadAll(res.Body)
	CheckError(err)

	var videoState ytVideoState
	json.Unmarshal(resBody, &videoState)

	if videoState.Error != "" {
		return videoState, errors.New(videoState.Error)
	}

	return videoState, nil
}

func formatCurrentSongInfo(state ytVideoState) string {
	if state.Player.TrackState == ytStatePaused {
		return "No song is currently playing"
	}

	videoUrl := ytVideoUrl + state.Video.Id
	if isSongCollection(state.Video.DurationSeconds) {
		videoUrl += ytTimeParam + strconv.Itoa(int(state.Player.VideoProgress))
	}

	if state.PlaylistId != "" {
		videoUrl += ytPlaylistParam + state.PlaylistId
	}

	return fmt.Sprintf("%s: %s %s", state.Video.Author, state.Video.Title, videoUrl)
}

func isSongCollection(durationSeconds int) bool {
	durationMinutes := SecondsToMinutes(durationSeconds)

	return durationMinutes >= songCollectionMinuteStart
}
