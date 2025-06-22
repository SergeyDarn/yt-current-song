package main

import (
	"flag"
	"fmt"

	"github.com/yt-current-song/pkg"
)

func main() {
	appIp := flag.String("appIp", "127.0.0.1", "YT Desktop app ip address")
	flag.Parse()

	fmt.Println(pkg.GetYtToken(*appIp))
}
