package internal

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
)

func BuildImage(log bool, no_cache bool) {
	var multiWriter io.Writer
	multiWriter = io.MultiWriter(os.Stdout) // 默认输出到控制台
	startTime := time.Now().Format("20060102-150405")

	if log {
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

	GenerateDockerfile(multiWriter) // 生成Dockerfile

	fmt.Fprintf(multiWriter, "=== 开始构建 [%s] ===\n", startTime)

	// 执行docker build命令
	cmdStr := "docker"
	args := []string{"build", "-t", "nightcord", "."}
	if no_cache {
		args = append(args, "--no-cache")
	}
	fmt.Fprint(multiWriter, "运行命令 ", cmdStr+strings.Join(args, " "))
	cmd := exec.Command(cmdStr, args...)
	cmd.Stdout = multiWriter
	cmd.Stderr = multiWriter

	// 记录构建开始时间
	buildStart := time.Now()
	// 使用buildStart来记录当前的构建开始时间
	fmt.Fprintf(multiWriter, "\n🚀 开始构建镜像 [%s]\n", buildStart.Format("2006-01-02 15:04:05"))
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(multiWriter, "❌ 构建失败: %v\n", err)
		return
	}
	fmt.Fprintln(multiWriter, "🎉 镜像构建完成")
	// 用buildStart计算实际构建耗时，并精确到三位小数
	duration := time.Since(buildStart)
	fmt.Fprintf(multiWriter, "⏱️ 构建耗时: %.3fs\n", duration.Seconds())
}
