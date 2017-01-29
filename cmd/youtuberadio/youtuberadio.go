package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/HugoGuiroux/youtuberadio/youtuberadio"
)

// var (
// 	message    = flag.String("message", "", "Text message to post")
// 	videoID    = flag.String("videoid", "", "ID of video to post")
// 	playlistID = flag.String("playlistid", "", "ID of playlist to post")
// )

func main() {
	flag.Parse()

	handle, err := youtuberadio.GetBroadcastUrl()
	if err != nil {
		log.Fatalf(err.Error())
		os.Exit(-1)
	}
	defer handle.Disconnect()

	info, err := handle.GetStreamInfo()
	if err != nil {
		log.Fatalf(err.Error())
		os.Exit(-1)
	}

	fmt.Println("Stream and broadcast bounded. URL is", info.URL, "and stream name is", info.Name)
}
