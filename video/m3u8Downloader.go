package VideoDownloader

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type SliceItem struct {
	Url   string
	Title string
}

var wg sync.WaitGroup
var cacheFolder = "tmp"
var downloadFolder = "DownloadedVideos"
var sliceItemDownloadChan = make(chan SliceItem, 3)
var count = 0
var progress = -1
var locker sync.Mutex

func init() {
	//创建缓存目录
	_, err := os.Stat(cacheFolder)
	if os.IsNotExist(err) {
		os.Mkdir(cacheFolder, 0777)
	} else {
		os.RemoveAll(cacheFolder)
		os.Mkdir(cacheFolder, 0777)
	}
	//创建存储目录
	_, err = os.Stat(downloadFolder)
	if os.IsNotExist(err) {
		os.Mkdir(downloadFolder, 0777)
	}
}

func M3u8FileDownloader(title string, playbackUrl string) {
	// title = strings.ReplaceAll(title, ":", "-")
	// m3u8 := m3u8_plugin.NewDownloader()
	// m3u8.SetUrl(playbackUrl)
	// //title = strings.ReplaceAll(title, " ", "")
	// m3u8.SetMovieName(title)
	// m3u8.SetNumOfThread(4)
	// m3u8.SetIfShowTheBar(true)
	// if m3u8.DefaultDownload() {
	// 	ffmpeg(title)
	// }
	downloadM3u8File(title, playbackUrl)
}

func downloadM3u8File(title string, playbackUrl string) {
	// 创建请求
	req, _ := http.NewRequest("GET", playbackUrl, nil)

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
	for key, value := range requestConfigs {
		req.Header.Set(key, value)
	}
	// 发送请求
	client := &http.Client{}
	resp, _ := client.Do(req)
	// 关闭响应
	defer resp.Body.Close()
	// 读取响应内容
	// body, _ := io.ReadAll(resp.Body)
	// result := string(body)
	// fmt.Println("mu3u8 file:", result)

	sliceUrl := parseSliceUrl(playbackUrl)

	fmt.Println("slice load url:", sliceUrl)
	if len(sliceUrl) == 0 {
		return
	}
	s := bufio.NewScanner(resp.Body)
	var lines []string
	for s.Scan() {
		lines = append(lines, s.Text())
	}

	var sliceTitles []string
	sliceItems := parseM3u8(lines, sliceUrl)
	count = len(sliceItems)
	addProgress()
	go downloadLoop()
	for index, item := range sliceItems {
		wg.Add(1)
		fileName := fmt.Sprintf("./%v/%v.ts", cacheFolder, strconv.Itoa(index))
		sliceTitles = append(sliceTitles, fileName)
		sliceItem := SliceItem{
			Url:   item,
			Title: fileName,
		}
		sliceItemDownloadChan <- sliceItem
	}
	wg.Wait()
	fmt.Println("\nall slice files downloaded")
	mergeTs := strings.Join(sliceTitles, "|")
	title = strings.ReplaceAll(title, ":", "-")
	mergeTsFiles(mergeTs, fmt.Sprintf("%v/%v.mp4", downloadFolder, title))
}

func parseSliceUrl(playbackUrl string) string {
	parsedURL, err := url.Parse(playbackUrl)
	if err != nil {
		log.Panic(err)
	}
	// 获取 scheme 和 host
	scheme := parsedURL.Scheme // "https"
	user := parsedURL.User     // dtliving-vip
	host := parsedURL.Host     // dingtalk.com
	path := parsedURL.Path     // "/live_hp/tt.m3u8"

	// 保留路径的基础部分
	basePath := path[:len("/live_hp/")]
	// 拼接成新的 URL
	if scheme == "http" {
		scheme = "https"
	}
	resultURL := fmt.Sprintf("%s://%s%s%s", scheme, user, host, basePath)
	return resultURL
}

func parseM3u8(lines []string, sliceUrl string) []string {
	var sliceItems []string
	fmt.Println(lines)
	for _, str := range lines {
		if str == "" {

		} else if strings.HasPrefix(str, "#") {

		} else {
			sliceItems = append(sliceItems, fmt.Sprintf("%v%v", sliceUrl, str))
		}
	}
	return sliceItems
}

func downloadLoop() {
	for {
		select {
		case item := <-sliceItemDownloadChan:
			downloadSliceItem(item.Title, item.Url)
		case <-time.After(10 * time.Minute):
			wg.Done()
		}
	}
}

func downloadSliceItem(fileName string, path string) {
	defer addProgress()
	defer wg.Done()
	resp, err := http.Get(path)
	if err != nil {
		log.Panic(err)
	}
	defer resp.Body.Close()
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Panic(err)
	}
	os.Remove(fileName)
	os.WriteFile(fileName, bytes, 0777)
}

func addProgress() {
	locker.Lock()
	progress = progress + 1
	printProgressBar(progress, count)
	locker.Unlock()
}
