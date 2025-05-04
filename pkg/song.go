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
	ytVideoTimeQuery          = "&t="
	ytPlaylistUrl             = "https://www.youtube.com/playlist?list="
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

	return formatCurrentSongInfo(songState.Video, songState.Player)
}

func GetCurrentPlaylistUrl(authToken string) string {
	songState, err := getYtVideoState(authToken)
	if err != nil {
		return err.Error()
	}

	return getPlaylistUrl(songState)
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

func formatCurrentSongInfo(video ytVideo, player ytPlayer) string {
	if player.TrackState == ytStatePaused {
		return "No song is currently playing"
	}

	videoUrl := ytVideoUrl + video.Id
	if isSongCollection(video.DurationSeconds) {
		videoUrl += ytVideoTimeQuery + strconv.Itoa(int(player.VideoProgress))
	}

	return fmt.Sprintf("%s: %s %s", video.Author, video.Title, videoUrl)
}

func getPlaylistUrl(videoState ytVideoState) string {
	if videoState.PlaylistId == "" {
		return "No playlist available"
	}

	if videoState.Player.TrackState == ytStatePaused {
		return "No song is currently playing"
	}

	return ytPlaylistUrl + videoState.PlaylistId
}

func isSongCollection(durationSeconds int) bool {
	durationMinutes := SecondsToMinutes(durationSeconds)

	return durationMinutes >= songCollectionMinuteStart
}
