package health

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func inspectSessionFile(path string) *SessionInfo {
	info := &SessionInfo{Path: path}

	stat, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			info.Exists = false
			return info
		}
		info.ParseError = err.Error()
		return info
	}

	info.Exists = true
	info.FileSizeBytes = stat.Size()
	info.UpdatedAt = stat.ModTime().Format(time.RFC3339)

	data, err := os.ReadFile(path)
	if err != nil {
		info.ParseError = err.Error()
		return info
	}

	messagesCount, updatedAt, err := parseSessionPayload(data)
	if err != nil {
		info.ParseError = err.Error()
		return info
	}

	info.MessagesCount = messagesCount
	if updatedAt != "" {
		info.UpdatedAt = updatedAt
	}
	return info
}

func inspectSessionsRoot(root string) *SessionsInfo {
	info := &SessionsInfo{Root: root}

	stat, err := os.Stat(root)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			info.Exists = false
			return info
		}
		info.ScanError = err.Error()
		return info
	}
	if !stat.IsDir() {
		info.ScanError = "sessions root is not a directory"
		return info
	}
	info.Exists = true

	walkErr := filepath.WalkDir(root, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			if info.ScanError == "" {
				info.ScanError = walkErr.Error()
			}
			return nil
		}
		if d.IsDir() || !strings.EqualFold(d.Name(), "session.json") {
			return nil
		}

		info.FilesCount++
		data, readErr := os.ReadFile(path)
		if readErr != nil {
			info.InvalidCount++
			info.InvalidFiles = append(info.InvalidFiles, SessionFileError{
				Path:       path,
				ParseError: readErr.Error(),
			})
			return nil
		}
		if _, _, parseErr := parseSessionPayload(data); parseErr != nil {
			info.InvalidCount++
			info.InvalidFiles = append(info.InvalidFiles, SessionFileError{
				Path:       path,
				ParseError: parseErr.Error(),
			})
			return nil
		}

		info.ValidCount++
		return nil
	})
	if walkErr != nil && info.ScanError == "" {
		info.ScanError = walkErr.Error()
	}

	return info
}

func parseSessionPayload(data []byte) (messagesCount int, updatedAt string, err error) {
	var payload struct {
		Messages  []json.RawMessage `json:"messages"`
		UpdatedAt string            `json:"updated_at"`
	}
	if err := json.Unmarshal(data, &payload); err != nil {
		return 0, "", err
	}
	return len(payload.Messages), strings.TrimSpace(payload.UpdatedAt), nil
}
