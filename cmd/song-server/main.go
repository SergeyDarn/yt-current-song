package main

import (
	"fmt"
	"net/http"

	yt "github.com/yt-current-song/pkg"
)

func getCurrentSongInfoRoute(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.NotFound(w, req)
		return
	}

	tokenQuery := req.URL.Query()["token"]
	if len(tokenQuery) == 0 {
		w.Write([]byte(yt.NoAuthTokenError))
		return
	}

	currentSongInfo := yt.GetCurrentSongInfo(tokenQuery[0])
	w.Write([]byte(currentSongInfo))
}

func main() {
	http.HandleFunc("/", getCurrentSongInfoRoute)

	err := http.ListenAndServe(":8050", nil)
	if err != nil {
		fmt.Printf("Couldn't start the server. Error: %s", err.Error())
	}
}
