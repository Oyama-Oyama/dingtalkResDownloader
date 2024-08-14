package VideoDownloader

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

func GetLiveRoomPublicInfo() {
	liveRoomUrl := "https://lv.dingtalk.com/getOpenLiveInfo?roomId=" + roomId + "&liveUuid=" + liveUuid

	// 创建请求
	req, _ := http.NewRequest("GET", liveRoomUrl, nil)

	byteData, _ := os.ReadFile("cookies.json")
	var cookies map[string]string
	_ = json.Unmarshal(byteData, &cookies)
	// 添加Cookies到请求
	var cookieStr strings.Builder
	for name, value := range cookies {
		cookieStr.WriteString(fmt.Sprintf("%s=%s; ", name, value))
	}
	cookieHeader := cookieStr.String()
	CookiepcSession := cookies["LV_PC_SESSION"]
	req.Header.Set("Cookie", cookieHeader)
	req.Header.Set("Cookie", "PC_SESSION="+CookiepcSession)
	// req.Header.Set("sec-ch-ua", "\"Not)A;Brand\";v=\"99\", \"Google Chrome\";v=\"127\", \"Chromium\";v=\"127\"")
	// req.Header.Set("sec-ch-ua-mobile", "?0")
	// req.Header.Set("sec-ch-ua-platform", "\"Windows\"")
	// req.Header.Set("sec-fetch-dest", "document")
	// req.Header.Set("sec-fetch-mode", "navigate")
	// req.Header.Set("sec-fetch-site", "none")

	// for key, value := range requestConfigs {
	// 	req.Header.Set(key, value)
	// }
	// 发送请求
	client := &http.Client{}
	resp, _ := client.Do(req)
	// 关闭响应
	defer resp.Body.Close()
	// 读取响应内容
	body, _ := io.ReadAll(resp.Body)

	var result map[string]interface{}
	json.Unmarshal(body, &result)

	//fmt.Println("live room info:", result)
	if len(result) == 0 {
		fmt.Println("error get live room info")
		return
	}
	title := result["openLiveDetailModel"].(map[string]interface{})["title"].(string)
	playbackUrl := result["openLiveDetailModel"].(map[string]interface{})["playbackUrl"].(string)

	// 获取当前时间并格式化
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	title += "_" + currentTime

	fmt.Println("标题:", title)
	fmt.Println("请求网址:", playbackUrl)
	M3u8FileDownloader(title, playbackUrl)
}
