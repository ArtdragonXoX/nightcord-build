package internal

import (
	"fmt"
	"io"
	"nightcord-build/utils"
)

func GetServerFile(tag string, w io.Writer) error {
	var err error
	if tag == "" {
		tag, err = utils.GetLatestReleaseTag("ArtDragonXoX", "nightcord-server")
		fmt.Fprintln(w, "获取最新版本号: ", tag)
		if err != nil {
			return err
		}
	} else {
		fmt.Fprintln(w, "指定版本号: ", tag)
	}
	fmt.Fprintln(w, "开始下载服务端文件")
	return utils.DownloadReleaseFiles("ArtDragonXoX", "nightcord-server", tag, []string{"nightcord-server"}, "./file")
}
