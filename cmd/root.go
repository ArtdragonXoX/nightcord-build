package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var (
	logFile   string
	no_cache  bool
	timeStamp bool
)

var rootCmd = &cobra.Command{
	Use:   "nightcord-build",
	Short: "Docker build pipeline manager",
	Long:  `Multi-stage Docker build system with logging and execution control`,
}

var makeCmd = &cobra.Command{
	Use:   "make",
	Short: "Generate Dockerfile",
	Long:  `Generate Dockerfile from template fragments`,
	Run: func(cmd *cobra.Command, args []string) {
		generateDockerfile(os.Stdout)
	},
}

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build Docker image",
	Long:  `Build Docker image with logging and cache optimization`,
	Run: func(cmd *cobra.Command, args []string) {
		startTime := time.Now().Format("20060102-150405")
		if logFile == "" {
			if timeStamp {
				logFile = fmt.Sprintf("build-%s.log", startTime)
			}
		}

		logWriter, err := os.Create(logFile)
		if err != nil {
			fmt.Printf("创建日志文件失败: %v\n", err)
			return
		}
		defer logWriter.Close()

		multiWriter := io.MultiWriter(os.Stdout, logWriter)
		fmt.Fprintf(multiWriter, "=== 开始构建 [%s] ===\n", startTime)
		buildImage(multiWriter)
	},
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run Docker container",
	Long:  `Run built Docker container with specified parameters`,
	Run: func(cmd *cobra.Command, args []string) {
		// 容器运行逻辑
	},
}

func generateDockerfile(w io.Writer) {
	dockerContent := &strings.Builder{}

	langFiles, err := filepath.Glob("langs/*.lang")
	if err != nil {
		fmt.Printf("查找.lang文件失败: %v\n", err)
		return
	}

	// 读取pre/post文件
	preContent, _ := os.ReadFile("dockerfile.pre")
	postContent, _ := os.ReadFile("dockerfile.post")

	// 多阶段构建模板
	dockerContent.WriteString("## 构建阶段\n")
	dockerContent.Write(preContent)

	for _, langFile := range langFiles {
		content, err := os.ReadFile(langFile)
		if err != nil {
			fmt.Printf("读取文件 %s 失败: %v\n", langFile, err)
			continue
		}

		dockerContent.WriteString(fmt.Sprintf("# ==== %s ====\n", langFile))
		dockerContent.Write(content)
		dockerContent.WriteString("\n\n")
	}

	dockerContent.WriteString("\n\n## 运行阶段\n")
	dockerContent.Write(postContent)

	// 写入Dockerfile
	if err := os.WriteFile("Dockerfile", []byte(dockerContent.String()), 0644); err != nil {
		fmt.Fprintf(w, "Dockerfile生成失败: %v\n", err)
		return
	}
	fmt.Fprintln(w, "✅ Dockerfile生成成功")
}

func buildImage(w io.Writer) {
	// 创建带时间戳的日志文件
	logFileHandle, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintf(w, "无法创建日志文件: %v\n", err)
		return
	}
	defer logFileHandle.Close()

	// 创建多路写入器
	multiWriter := io.MultiWriter(w, logFileHandle)

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

	startTime := time.Now().Format("2006-01-02 15:04:05")
	fmt.Fprintf(multiWriter, "\n🚀 开始构建镜像 [%s]\n", startTime)
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(multiWriter, "❌ 构建失败: %v\n", err)
		return
	}
	fmt.Fprintln(multiWriter, "🎉 镜像构建完成")
	fmt.Fprintf(multiWriter, "⏱️ 构建耗时: %s\n", time.Since(time.Now()).Round(time.Second))
}

func init() {
	rootCmd.AddCommand(makeCmd)
	rootCmd.AddCommand(buildCmd)
	rootCmd.AddCommand(runCmd)

	buildCmd.Flags().StringVarP(&logFile, "log", "l", "build.log", "日志文件路径 (默认: ./build.log)")
	buildCmd.Flags().BoolVarP(&timeStamp, "timestamp", "t", false, "在日志文件名中添加时间戳")
	buildCmd.Flags().BoolVarP(&no_cache, "no-cache", "n", false, "不使用缓存构建镜像")

}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
