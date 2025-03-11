package cmd

import (
	"fmt"
	"nightcord-build/internal"
	"nightcord-build/internal/model"
	"os"

	"github.com/spf13/cobra"
)

var (
	conf = model.Config{}
)

var rootCmd = &cobra.Command{
	Use:   "nightcord-build",
	Short: "Docker build pipeline manager",
	Long:  `Multi-stage Docker build system with logging and execution control`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if conf.Log {
			internal.LogEnable()
		}
	},
}

var makeCmd = &cobra.Command{
	Use:   "make",
	Short: "Generate Dockerfile",
	Long:  `Generate Dockerfile from template fragments`,
	Run: func(cmd *cobra.Command, args []string) {
		internal.GenerateDockerfile(conf)
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
		internal.Run(conf)
	},
}

func init() {
	rootCmd.AddCommand(makeCmd)
	rootCmd.AddCommand(buildCmd)
	rootCmd.AddCommand(runCmd)

	makeCmd.Flags().BoolVarP(&conf.Log, "log", "l", false, "生成日志文件")
	makeCmd.Flags().BoolVarP(&conf.Dev, "dev", "d", false, "搭建开发环境")
	makeCmd.Flags().StringVarP(&conf.Repo, "repo", "r", "", "仓库地址")

	buildCmd.Flags().BoolVarP(&conf.Log, "log", "l", false, "生成日志文件")
	buildCmd.Flags().BoolVarP(&conf.NoCache, "no-cache", "n", false, "不使用缓存构建镜像")
	buildCmd.Flags().StringVarP(&conf.Tag, "tag", "t", "", "服务端标签")
	buildCmd.Flags().BoolVarP(&conf.LocalFile, "local", "f", false, "使用本地服务端文件")
	buildCmd.Flags().StringVarP(&conf.LocalFilePath, "local-file", "p", "", "本地服务端文件路径")
	buildCmd.Flags().BoolVarP(&conf.Dev, "dev", "d", false, "搭建开发环境")
	buildCmd.Flags().StringVarP(&conf.Repo, "repo", "r", "", "仓库地址")

	runCmd.Flags().BoolVarP(&conf.Log, "log", "l", false, "生成日志文件")
	runCmd.Flags().BoolVarP(&conf.NoCache, "no-cache", "n", false, "不使用缓存构建镜像")
	runCmd.Flags().StringVarP(&conf.Tag, "tag", "t", "", "服务端标签")
	runCmd.Flags().BoolVarP(&conf.LocalFile, "local", "f", false, "使用本地服务端文件")
	runCmd.Flags().StringVarP(&conf.LocalFilePath, "local-file", "p", "", "本地服务端文件路径")
	runCmd.Flags().BoolVarP(&conf.Dev, "dev", "d", false, "搭建开发环境")
	runCmd.Flags().StringVarP(&conf.Repo, "repo", "r", "", "仓库地址")
	runCmd.Flags().StringVarP(&conf.Volume, "volume", "v", "", "挂载卷")
	runCmd.Flags().BoolVarP(&conf.JumpBuild, "jump-build", "j", false, "跳过构建镜像步骤")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
