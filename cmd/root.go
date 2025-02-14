package cmd

import (
	"fmt"
	"nightcord-build/internal"
	"os"

	"github.com/spf13/cobra"
)

var (
	conf = internal.Config{}
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
		internal.GenerateDockerfile(os.Stdout)
	},
}

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build Docker image",
	Long:  `Build Docker image with logging and cache optimization`,
	Run: func(cmd *cobra.Command, args []string) {
		internal.BuildImage(conf)
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

func init() {
	rootCmd.AddCommand(makeCmd)
	rootCmd.AddCommand(buildCmd)
	rootCmd.AddCommand(runCmd)

	buildCmd.Flags().BoolVarP(&conf.Log, "log", "l", false, "生成日志文件")
	buildCmd.Flags().BoolVarP(&conf.NoCache, "no-cache", "n", false, "不使用缓存构建镜像")
	buildCmd.Flags().StringVarP(&conf.Tag, "tag", "t", "", "服务端标签")
	buildCmd.Flags().BoolVarP(&conf.LocalFile, "local", "f", false, "使用本地服务端文件")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
