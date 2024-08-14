**下载钉钉视频**
[确保已安装Chrome浏览器]
***1. 获取视频分享链接***
![获取分享链接](pic/video_share.png)
***2. 下载视频***
在main目录下执行
``` go run main.go -type "video" -url "https://n.dingtalk.com/dingding/live-room/index.html?roomId=oiPzbfYe1oANm42n&liveUuid=09bd1176-fff2-4a51-9677-1ecf436ebe48" ```


***golang 依赖下载失败时可以修改为国内镜像再次尝试:***
```setx GOPROXY "https://mirrors.aliyun.com/goproxy"```