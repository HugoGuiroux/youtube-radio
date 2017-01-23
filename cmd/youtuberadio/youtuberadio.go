package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"google.golang.org/api/youtube/v3"
	"github.com/HugoGuiroux/youtuberadio/youtuberadio"
)

// var (
// 	message    = flag.String("message", "", "Text message to post")
// 	videoID    = flag.String("videoid", "", "ID of video to post")
// 	playlistID = flag.String("playlistid", "", "ID of playlist to post")
// )

func main() {
	flag.Parse()

	client, err := youtuberadio.BuildOAuthHTTPClient(youtube.YoutubeScope)
	if err != nil {
		log.Fatalf("Error building OAuth client: %v", err)
	}

	service, err := youtube.New(client)
	if err != nil {
		log.Fatalf("Error creating YouTube client: %v", err)
	}

	broadcastRequest := &youtube.LiveBroadcast{
		Snippet: &youtube.LiveBroadcastSnippet{
			Title: "Test LiveBroadcast",
			ScheduledStartTime: time.Now().Format("2006-01-02T15:04:05-0700"),
			ScheduledEndTime: time.Now().Add(time.Hour).Format("2006-01-02T15:04:05-0700"),
		},
		Status: &youtube.LiveBroadcastStatus{
			PrivacyStatus: "private",
		},
	}
	callBroadcast := service.LiveBroadcasts.Insert("snippet,status", broadcastRequest)
	broadcast, err := callBroadcast.Do()
	if err != nil {
		log.Fatalf("Error making API call to create broadcast: %v", err.Error())
	}

	bid := broadcast.Id
	defer func() {
		fmt.Println("Deleting broadcast", bid)
		service.LiveBroadcasts.Delete(bid).Do()
	}()
	fmt.Printf("Your broadcast ID is %v\n", bid)

	streamRequest := &youtube.LiveStream{
		Snippet: &youtube.LiveStreamSnippet{
			Title: "Test LiveStream",
		},
		Cdn: &youtube.CdnSettings{
			FrameRate: "30fps",
			Resolution: "240p",
			IngestionType: "rtmp",
		},
	}
	callStream := service.LiveStreams.Insert("snippet,cdn", streamRequest)
	stream, err := callStream.Do()
	if err != nil {
		log.Fatalf("Error making API call to create stream: %v", err.Error())
	}

	sid := stream.Id
	defer func() {
		fmt.Println("Deleting stream", sid)
		service.LiveStreams.Delete(sid).Do()
	}()
	fmt.Printf("Your stream ID is %v\n", sid)

	callBind := service.LiveBroadcasts.Bind(bid, "id")
	callBind.StreamId(sid)
	_, err = callBind.Do()
	if err != nil {
		log.Fatalf("Error making API call to bind broadcast and stream: %v", err.Error())
	}

	fmt.Println("Stream and broadcast bounded. URL is", stream.Cdn.IngestionInfo.IngestionAddress, "and stream name is", stream.Cdn.IngestionInfo.StreamName)
}
