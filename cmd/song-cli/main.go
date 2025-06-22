package main

import (
	"flag"
	"fmt"

	yt "github.com/yt-current-song/pkg"
)

func main() {
	authToken := flag.String("token", "", "API Token to connect to YT Desktop app")
	appIp := flag.String("appIp", "127.0.0.1", "YT Desktop app ip address")
	flag.Parse()

	if *authToken == "" {
		fmt.Println(yt.NoAuthTokenError)
		return
	}

	fmt.Println(yt.GetCurrentSongInfo(*authToken, *appIp))
}
