package utils_test

import (
	"nightcord-build/utils"
	"testing"
)

func TestCloneRepo(t *testing.T) {
	repo := "https://github.com/ArtdragonXoX/nightcord-server.git"
	destDir := "testdata/nightcord-server"
	err := utils.CloneRepo(repo, destDir)
	if err != nil {
		t.Errorf("CloneRepo() error = %v", err)
	}
}
