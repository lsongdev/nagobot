package health

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"time"

	"github.com/linanwx/nagobot/cron"
)

func inspectCronFile(path string) *CronInfo {
	info := &CronInfo{Path: path}

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

	f, err := os.Open(path)
	if err != nil {
		info.ParseError = err.Error()
		return info
	}
	defer f.Close()

	var jobs []cron.Job
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := bytes.TrimSpace(scanner.Bytes())
		if len(line) == 0 {
			continue
		}
		var job cron.Job
		if err := json.Unmarshal(line, &job); err != nil {
			info.ParseError = err.Error()
			return info
		}
		jobs = append(jobs, job)
	}
	if err := scanner.Err(); err != nil {
		info.ParseError = err.Error()
		return info
	}

	info.JobsCount = len(jobs)
	return info
}
