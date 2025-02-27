package internal

import (
	"fmt"
	"nightcord-build/internal/model"
	"nightcord-build/utils"
	"os"
	"os/exec"
	"strings"
	"time"
)

func BuildImage(conf model.Config) {
	startTime := time.Now().Format("20060102-150405")

	if conf.Dev {
		fmt.Fprintln(multiWriter, "构建nightcord-server开发环境镜像")
	} else {
		fmt.Fprintln(multiWriter, "🚀 开始获取服务端文件")
		if conf.LocalFilePath != "" {
			fmt.Fprintln(multiWriter, "🔍 使用本地服务端文件")
			fmt.Fprintln(multiWriter, "📂 复制本地服务端文件")
			if err := os.MkdirAll("file", 0755); err != nil {
				fmt.Fprintf(multiWriter, "创建文件目录失败: %v\n", err)
				panic(err)
			}
			_ = os.Remove("file/nightcord-server") // 先尝试删除已有文件（忽略错误）
			if err := utils.CopyFile(conf.LocalFilePath, "file/nightcord-server"); err != nil {
				fmt.Fprintf(multiWriter, "❌ 复制本地服务端文件失败: %v\n", err)
				panic(err)
			}
			fmt.Fprintln(multiWriter, "🎉 复制本地服务端文件成功")
		} else {
			// 检查本地文件模式
			if conf.LocalFile {
				if _, err := os.Stat("file/nightcord-server"); os.IsNotExist(err) {
					fmt.Fprintf(multiWriter, "❌ 本地服务端文件不存在: file/nightcord-server\n")
					panic(err)
				}
				fmt.Fprintln(multiWriter, "🔍 使用本地服务端文件")
			} else {
				fmt.Fprintln(multiWriter, "🌐 从GitHub获取服务端文件")
				err := GetServerFile(conf.Tag, multiWriter)
				if err != nil {
					fmt.Fprintf(multiWriter, "❌ 获取服务端文件失败: %v\n", err)
					panic(err)
				}
				fmt.Fprintln(multiWriter, "🎉 获取服务端文件成功")
			}
		}
	}
	GenerateDockerfile(conf) // 生成Dockerfile

	fmt.Fprintf(multiWriter, "=== 开始构建 [%s] ===\n", startTime)

	// 执行docker build命令
	cmdStr := "docker"
	var containerName string
	if conf.Dev {
		containerName = "nightcord-dev"
	} else {
		containerName = "nightcord"
	}
	args := []string{"build", "-t", containerName, "."}
	if conf.NoCache {
		args = append(args, "--no-cache")
	}
	fmt.Fprintln(multiWriter, "运行命令 ", cmdStr+" "+strings.Join(args, " "))
	cmd := exec.Command(cmdStr, args...)
	cmd.Stdout = multiWriter
	cmd.Stderr = multiWriter

	// 记录构建开始时间
	buildStart := time.Now()
	// 使用buildStart来记录当前的构建开始时间
	fmt.Fprintf(multiWriter, "\n🚀 开始构建镜像 [%s]\n", buildStart.Format("2006-01-02 15:04:05"))
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(multiWriter, "❌ 构建失败: %v\n", err)
		panic(err)
	}
	fmt.Fprintln(multiWriter, "🎉 镜像构建完成")
	// 用buildStart计算实际构建耗时，并精确到三位小数
	duration := time.Since(buildStart)
	fmt.Fprintf(multiWriter, "⏱️ 构建耗时: %.3fs\n", duration.Seconds())
}
