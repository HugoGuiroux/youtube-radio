package youtuberadio

import (
	"fmt"
	"net/http"
	"time"

	"google.golang.org/api/youtube/v3"
)

type Handle struct {
	client    *http.Client
	service   *youtube.Service
	broadcast *youtube.LiveBroadcast
	stream    *youtube.LiveStream
}

func GetBroadcastUrl() (*Handle, error) {
	var handle = &Handle{}
	var err error
	handle.client, err = BuildOAuthHTTPClient(youtube.YoutubeScope)

	if err != nil {
		return nil, fmt.Errorf("Error building OAuth client: %v", err.Error())
	}

	handle.service, err = youtube.New(handle.client)
	if err != nil {
		return nil, fmt.Errorf("Error creating YouTube client: %v", err.Error())
	}

	broadcastRequest := &youtube.LiveBroadcast{
		Snippet: &youtube.LiveBroadcastSnippet{
			Title:              "Test LiveBroadcast",
			ScheduledStartTime: time.Now().Format("2006-01-02T15:04:05-0700"),
			ScheduledEndTime:   time.Now().Add(time.Hour).Format("2006-01-02T15:04:05-0700"),
		},
		Status: &youtube.LiveBroadcastStatus{
			PrivacyStatus: "private",
		},
	}
	callBroadcast := handle.service.LiveBroadcasts.Insert("snippet,status", broadcastRequest)
	handle.broadcast, err = callBroadcast.Do()
	if err != nil {
		return nil, fmt.Errorf("Error making API call to create broadcast: %v", err.Error())
	}

	bid := handle.broadcast.Id
	streamRequest := &youtube.LiveStream{
		Snippet: &youtube.LiveStreamSnippet{
			Title: "Test LiveStream",
		},
		Cdn: &youtube.CdnSettings{
			FrameRate:     "30fps",
			Resolution:    "240p",
			IngestionType: "rtmp",
		},
	}
	callStream := handle.service.LiveStreams.Insert("snippet,cdn", streamRequest)
	handle.stream, err = callStream.Do()
	if err != nil {
		handle.service.LiveBroadcasts.Delete(bid).Do()
		return nil, fmt.Errorf("Error making API call to create stream: %v", err.Error())
	}

	sid := handle.stream.Id

	callBind := handle.service.LiveBroadcasts.Bind(bid, "id")
	callBind.StreamId(sid)
	_, err = callBind.Do()
	if err != nil {
		handle.service.LiveStreams.Delete(sid).Do()
		return nil, fmt.Errorf("Error making API call to bind broadcast and stream: %v", err.Error())
	}

	return handle, nil
}

func (h *Handle) Disconnect() {
	if h.service != nil && h.stream != nil {
		h.service.LiveStreams.Delete(h.stream.Id).Do()
	}

	if h.service != nil && h.broadcast != nil {
		h.service.LiveBroadcasts.Delete(h.broadcast.Id).Do()
	}
}

type StreamInfo struct {
	URL  string
	Name string
}

func (h *Handle) GetStreamInfo() (*StreamInfo, error) {
	if h.stream == nil {
		return nil, fmt.Errorf("Not connected to the stream")
	}
	return &StreamInfo{
		URL:  h.stream.Cdn.IngestionInfo.IngestionAddress,
		Name: h.stream.Cdn.IngestionInfo.StreamName,
	}, nil
}
