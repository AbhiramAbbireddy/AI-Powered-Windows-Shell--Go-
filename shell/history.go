package shell

import (
	"fmt"
	"os"
	"path/filepath"
)

type HistoryEntry struct {
	UserInput string
	Command   string
}

func AppendHistory(historyPath string, entry HistoryEntry) error {
	if err := os.MkdirAll(filepath.Dir(historyPath), 0o755); err != nil {
		return err
	}

	file, err := os.OpenFile(historyPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()

	record := fmt.Sprintf("input=%q | command=%q\n", entry.UserInput, entry.Command)
	_, err = file.WriteString(record)
	return err
}
