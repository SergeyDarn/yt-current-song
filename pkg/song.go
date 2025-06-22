package pkg

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

const (
	NoAuthTokenError = "Cannot access current song info without YT Desktop api token"
	authHeader       = "authorization"

	ytDesktopApiUrl      = "http://{ip}:9863/api/v1/"
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

func (s ytVideoState) getFormattedSongInfo() string {
	if s.Player.TrackState == ytStatePaused {
		return "No song is currently playing"
	}

	return fmt.Sprintf(
		"%s: %s %s",
		s.Video.Author,
		s.Video.Title,
		s.getVideoUrl(),
	)
}

func (s ytVideoState) getVideoUrl() string {
	videoUrl := ytVideoUrl + s.Video.Id

	s.addVideoTime(&videoUrl)
	s.addVideoPlaylist(&videoUrl)

	return videoUrl
}

func (s ytVideoState) addVideoTime(videoUrl *string) {
	videoDurationMinutes := SecondsToMinutes(s.Video.DurationSeconds)
	isSongCollection := videoDurationMinutes >= songCollectionMinuteStart

	if isSongCollection {
		videoTime := int(s.Player.VideoProgress)
		*videoUrl += ytTimeParam + strconv.Itoa(videoTime)
	}
}

func (s ytVideoState) addVideoPlaylist(videoUrl *string) {
	if s.PlaylistId != "" {
		*videoUrl += ytPlaylistParam + s.PlaylistId
	}
}

func GetCurrentSongInfo(authToken string, appIp string) string {
	songState, err := getYtVideoState(authToken, appIp)
	if err != nil {
		return err.Error()
	}

	return songState.getFormattedSongInfo()
}

func getYtVideoState(authToken string, appIp string) (ytVideoState, error) {
	url := strings.Replace(ytDesktopGetStateUrl, "{ip}", appIp, 1)
	req, err := http.NewRequest("GET", url, nil)
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
