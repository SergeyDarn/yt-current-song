package main

import (
	"flag"
	"fmt"

	"github.com/yt-current-song/pkg"
)

func main() {
	authToken := flag.String("token", "", "API Token to connect to YT Desktop app")
	flag.Parse()

	if *authToken == "" {
		fmt.Println(pkg.NoAuthTokenError)
		return
	}

	fmt.Println(pkg.GetCurrentPlaylistUrl(*authToken))
}
