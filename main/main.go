package main

import (
	"DingTalk/VideoDownloader"
	"flag"
	"fmt"
)

func main() {
	fileType := flag.String("type", "", "aaa")
	url := flag.String("url", "", "链接")
	flag.Parse()
	fmt.Printf("%v with url:%v start to parse\n", *fileType, *url)
	if *fileType == "video" {
		VideoDownloader.StartWithUrl(*url)
	} else {
		fmt.Println("invalid file type:%v", *fileType)
	}
}
