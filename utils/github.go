package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// GetLatestReleaseTag returns the latest release tag for a GitHub repository
func GetLatestReleaseTag(owner, repo string) (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", owner, repo)

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch release: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var release struct {
		TagName string `json:"tag_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if release.TagName == "" {
		return "", errors.New("no tag name found in release")
	}

	return release.TagName, nil
}

// DownloadReleaseFiles downloads specified files from a GitHub release
func DownloadReleaseFiles(owner, repo, tag string, filenames []string, destDir string) error {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/tags/%s", owner, repo, tag)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch release: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var release struct {
		Assets []struct {
			Name string `json:"name"`
			URL  string `json:"browser_download_url"`
		} `json:"assets"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	// Create destination directory if it doesn't exist
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Download each requested file
	for _, filename := range filenames {
		var downloadURL string
		for _, asset := range release.Assets {
			if asset.Name == filename {
				downloadURL = asset.URL
				break
			}
		}

		if downloadURL == "" {
			return fmt.Errorf("file %s not found in release", filename)
		}

		resp, err := http.Get(downloadURL)
		if err != nil {
			return fmt.Errorf("failed to download file %s: %w", filename, err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("failed to download file %s: status code %d", filename, resp.StatusCode)
		}

		filePath := filepath.Join(destDir, filename)
		out, err := os.Create(filePath)
		if err != nil {
			return fmt.Errorf("failed to create file %s: %w", filePath, err)
		}
		defer out.Close()

		if _, err := io.Copy(out, resp.Body); err != nil {
			return fmt.Errorf("failed to write file %s: %w", filePath, err)
		}
	}

	return nil
}
