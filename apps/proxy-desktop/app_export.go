package main

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// SaveExportFile opens a save dialog and writes base64-encoded file contents.
// Returns the chosen path, or empty string if the user cancelled.
func (a *App) SaveExportFile(defaultFilename string, contentBase64 string) (string, error) {
	if a.ctx == nil {
		return "", fmt.Errorf("application not ready")
	}
	path, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		DefaultFilename: defaultFilename,
	})
	if err != nil {
		return "", err
	}
	if path == "" {
		return "", nil
	}
	data, err := base64.StdEncoding.DecodeString(contentBase64)
	if err != nil {
		return "", fmt.Errorf("decode export data: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return "", fmt.Errorf("write export file: %w", err)
	}
	return path, nil
}

// PickDatabaseFile opens a file dialog to choose a SQLite database file.
func (a *App) PickDatabaseFile() (string, error) {
	if a.ctx == nil {
		return "", fmt.Errorf("application not ready")
	}
	return runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select proxy database",
		Filters: []runtime.FileFilter{
			{DisplayName: "SQLite database", Pattern: "*.db"},
			{DisplayName: "All files", Pattern: "*"},
		},
	})
}
