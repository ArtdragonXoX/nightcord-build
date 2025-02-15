package internal

import (
	"fmt"
	"os"
	"os/exec"
)

func Run(conf Config) {
	if _, err := os.Stat("docker-compose.yaml"); os.IsNotExist(err) {
		fmt.Fprintf(multiWriter, "❌ docker-compose.yaml不存在，请先创建docker-compose.yaml文件\n")
		return
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
		return
	}
	fmt.Fprintf(multiWriter, "🎉 容器启动成功\n")
}
