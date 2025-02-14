package internal

import "nightcord-build/utils"

func GetServerFile(tag string) error {
	var err error
	if tag == "" {
		tag, err = utils.GetLatestReleaseTag("ArtDragonXoX", "nightcord-server")
		if err != nil {
			return err
		}
	}
	return utils.DownloadReleaseFiles("ArtDragonXoX", "nightcord-server", tag, []string{"nightcord-server"}, "./file")
}
