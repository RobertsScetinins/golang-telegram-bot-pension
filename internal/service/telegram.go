package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type FileInfo struct {
	OK     bool `json:ok`
	Result models.File
}

func GetFileInfo(url string) (*FileInfo, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	var fileInfo FileInfo
	if err := json.Unmarshal(body, &fileInfo); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if !fileInfo.OK {
		return nil, fmt.Errorf("telegram API returned error")
	}

	return &fileInfo, nil
}

func GetDownloadLink(b *bot.Bot, url string) (string, error) {
	fileInfo, err := GetFileInfo(url)
	if err != nil {
		return "", err
	}

	downloadLink := b.FileDownloadLink(&fileInfo.Result)

	return downloadLink, nil
}
