package infrastructure

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type TransmissionDownloader struct{}

func NewTransmissionDownloader() *TransmissionDownloader {
	return &TransmissionDownloader{}
}

func (d *TransmissionDownloader) Download(filePath string) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", filePath)
	}

	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	cmd := exec.Command("transmission-cli", absPath)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to download torrent: %w, output: %s", err, output)
	}

	return nil
}
