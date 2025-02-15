package internal

import (
	"fmt"
	"io"
	"os"
	"time"
)

var multiWriter = io.MultiWriter(os.Stdout)

func LogEnable() {
	startTime := time.Now().Format("20060102-150405")
	logFile := fmt.Sprintf("./logs/build-%s.log", startTime)
	if err := os.MkdirAll("./logs", 0755); err != nil {
		fmt.Printf("无法创建日志目录: %v\n", err)
		return
	}

	logWriter, err := os.Create(logFile)
	if err != nil {
		fmt.Printf("创建日志文件失败: %v\n", err)
		return
	}
	defer logWriter.Close()

	multiWriter = io.MultiWriter(os.Stdout, logWriter) // 同时输出到控制台和日志文件
}
