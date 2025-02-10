package internal

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func GenerateDockerfile(w io.Writer) {
	fmt.Fprintln(w, "🔍 开始生成Dockerfile")
	dockerContent := &strings.Builder{}

	langFiles, err := filepath.Glob("langs/*.lang")
	if err != nil {
		fmt.Fprintf(w, "❌ 查找.lang文件失败: %v\n", err)
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
			fmt.Fprintf(w, "❌ 读取文件 %s 失败: %v\n", langFile, err)
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
		fmt.Fprintf(w, "❌ Dockerfile生成失败: %v\n", err)
		return
	}
	fmt.Fprintln(w, "✅ Dockerfile生成成功")
}
