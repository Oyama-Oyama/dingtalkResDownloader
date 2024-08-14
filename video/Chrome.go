package VideoDownloader

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

func StartChromeWithUrl(path string) {
	fmt.Printf("start with url:%v\n", path)

	opts := append(
		chromedp.DefaultExecAllocatorOptions[3:],
		chromedp.NoDefaultBrowserCheck,
		chromedp.NoFirstRun,
	)

	parentCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(parentCtx)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()

	var siteCookies []*network.Cookie
	//	H5url := "https://h5.dingtalk.com"
	chromedp.Run(
		ctx,
		network.Enable(),
		chromedp.Navigate(path),
		chromedp.WaitVisible("video", chromedp.ByQueryAll),
		chromedp.ActionFunc(func(ctx context.Context) error {
			siteCookies, _ = network.GetCookies().Do(ctx)
			fmt.Println("cookies loaded:", siteCookies)
			return nil
		}),
	)
	os.Remove("cookies.json")
	// 保存cookies到文件
	cookies := make(map[string]string)
	for _, cookie := range siteCookies {
		cookies[cookie.Name] = cookie.Value
	}
	jsonCookies, _ := json.Marshal(cookies)
	os.WriteFile("cookies.json", jsonCookies, 0644)
}
