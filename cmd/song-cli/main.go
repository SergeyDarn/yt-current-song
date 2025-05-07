package main

import (
	"flag"
	"fmt"

	yt "github.com/yt-current-song/pkg"
)

func main() {
	authToken := flag.String("token", "", "API Token to connect to YT Desktop app")
	flag.Parse()

	if *authToken == "" {
		fmt.Println(yt.NoAuthTokenError)
		return
	}

	fmt.Println(yt.GetCurrentSongInfo(*authToken))
}
