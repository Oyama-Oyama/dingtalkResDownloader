package VideoDownloader

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
)

var roomId string
var liveUuid string
var requestConfigs map[string]string

func init() {
	byteData, _ := os.ReadFile("config.json")
	_ = json.Unmarshal(byteData, &requestConfigs)
}

func parseUrl(path string) {
	parsedURL, err := url.Parse(path)
	if err != nil {
		fmt.Println("parse url failed:", err)
		return
	}
	queryParams := parsedURL.Query()
	roomId = queryParams.Get("roomId")
	liveUuid = queryParams.Get("liveUuid")
	if liveUuid == "" || roomId == "" {
		fmt.Printf("error liveUuid=%v  roomId=%v\v", liveUuid, roomId)
		return
	}
}

func StartWithUrl(url string) {
	StartChromeWithUrl(url)
	parseUrl(url)
	GetLiveRoomPublicInfo()
}
