package main

import (
	"fmt"
	"net/http"

	"github.com/yt-current-song/pkg"
)

func getCurrentSongInfoRoute(w http.ResponseWriter, req *http.Request) {
	tokenQuery := req.URL.Query()["token"]

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	if len(tokenQuery) == 0 {
		w.Write([]byte("Не могу получить информацию о текущей песне без токена авторизации"))
		return
	}

	currentSongInfo := pkg.GetCurrentSongInfo(tokenQuery[0])
	w.Write([]byte(currentSongInfo))
}

func main() {
	http.HandleFunc("/", getCurrentSongInfoRoute)

	err := http.ListenAndServe(":8050", nil)
	if err != nil {
		fmt.Printf("Couldn't start the server. Error: %s", err.Error())
	}
}
