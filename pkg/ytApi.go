package pkg

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

const (
	appId = "yt-current-song"

	ytApiUrl      = "http://{ip}:9863/api/v1/"
	ytGetStateUrl = ytApiUrl + "state"

	ytRequestCodeUrl  = ytApiUrl + "auth/requestcode"
	ytRequestTokenUrl = ytApiUrl + "auth/request"
)

type ytCodeResponse struct {
	Code  string
	Error string
}

type ytTokenResponse struct {
	Token string
	Error string
}

func GetYtVideoState(authToken string, appIp string) (ytVideoState, error) {
	url := replaceAppIp(ytGetStateUrl, appIp)
	headers := map[string]string{"authorization": authToken}

	resBody := doRequest(http.MethodGet, url, nil, headers)

	var videoState ytVideoState
	json.Unmarshal(resBody, &videoState)

	if videoState.Error != "" {
		return videoState, errors.New(videoState.Error)
	}

	return videoState, nil
}

func GetYtToken(appIp string) string {
	ytCode := getYtCode(appIp)

	body := fmt.Appendf([]byte{}, `{
		"appId": "%s",
		"code": "%s"
	}`, appId, ytCode)
	url := replaceAppIp(ytRequestTokenUrl, appIp)
	resBody := doRequest(http.MethodPost, url, body, map[string]string{})

	var ytToken ytTokenResponse
	json.Unmarshal(resBody, &ytToken)

	if ytToken.Error != "" {
		panic(ytToken.Error)
	}

	return ytToken.Token
}

func getYtCode(appIp string) string {
	body := fmt.Appendf([]byte{}, `{
		"appId": "%s",
		"appName": "YT Current Song",
		"appVersion": "1.0.0"
	}`, appId)

	url := replaceAppIp(ytRequestCodeUrl, appIp)
	resBody := doRequest(http.MethodPost, url, body, map[string]string{})

	var ytCode ytCodeResponse
	json.Unmarshal(resBody, &ytCode)

	if ytCode.Error != "" {
		panic(ytCode.Error)
	}

	return ytCode.Code
}

func replaceAppIp(url string, appIp string) string {
	return strings.Replace(url, "{ip}", appIp, 1)
}
