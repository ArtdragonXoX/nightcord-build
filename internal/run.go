package internal

import (
	"fmt"
	"nightcord-build/internal/model"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func Run(conf model.Config) {
	if conf.Dev {
		fmt.Fprintln(multiWriter, "构建nightcord-server开发环境容器")
		BuildImage(conf)
		fmt.Fprintln(multiWriter, "🚀 正在启动容器...")
		cmdStr := "docker"
		args := []string{"run", "--name", "nightcord-dev"}

		// 保留交互式终端参数
		if runtime.GOOS != "windows" {
			args = append(args, "-it")
		} else {
			args = append(args, "-i") // Windows下只保留输入参数
		}

		if conf.Volume != "" {
			args = append(args, "-v", fmt.Sprintf("%s:/home/nightcord-server", conf.Volume))
		}
		args = append(args, "nightcord-dev:latest")

		// Windows特殊处理
		if runtime.GOOS == "windows" {
			if _, err := exec.LookPath("winpty"); err == nil {
				cmdStr = "winpty"
				args = append([]string{"docker"}, args...)
			} else {
				fmt.Fprintln(multiWriter, "⚠️  检测到Windows环境但未找到winpty，使用基础模式运行")
			}
		}

		fmt.Fprintln(multiWriter, "运行命令 ", cmdStr+" "+strings.Join(args, " "))
		cmd := exec.Command(cmdStr, args...)
		cmd.Stdout = multiWriter
		cmd.Stderr = multiWriter
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(multiWriter, "❌ 启动容器失败: %v\n", err)
			panic(err)
		}
	} else {
		if _, err := os.Stat("docker-compose.yaml"); os.IsNotExist(err) {
			fmt.Fprintf(multiWriter, "❌ docker-compose.yaml不存在，请先创建docker-compose.yaml文件\n")
			panic(err)
		}
		BuildImage(conf)
		cmdStr := "docker-compose"
		args := []string{"up", "-d"}
		fmt.Fprintf(multiWriter, "🚀 正在启动容器...\n")
		cmd := exec.Command(cmdStr, args...)
		cmd.Stdout = multiWriter
		cmd.Stderr = multiWriter
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(multiWriter, "❌ 启动容器失败: %v\n", err)
			panic(err)
		}
	}
	fmt.Fprintf(multiWriter, "🎉 容器启动成功\n")
}
