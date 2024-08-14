package VideoDownloader

import "fmt"

func printProgressBar(completed int, total int) {
	barLength := 50                                 // 定义进度条的长度
	progress := float64(completed) / float64(total) // 计算进度
	filled := int(progress * float64(barLength))    // 已完成的部分长度
	empty := barLength - filled                     // 剩余部分长度

	// 输出进度条
	fmt.Printf("\r[")
	for i := 0; i < filled; i++ {
		fmt.Print("=") // 已完成部分
	}
	for i := 0; i < empty; i++ {
		fmt.Print(" ") // 剩余部分
	}
	fmt.Printf("] %.2f%%", progress*100) // 输出百分比
	// Flush the buffer
	//fmt.Flush()
}
