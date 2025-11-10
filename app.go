package main

import (
	"context"
	"cropper/core"      //  models
	"cropper/processor" //  video engine
	"cropper/utils"     //  helpers
	"fmt"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) OpenFileDialog() (string, error) {
	filePath, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select Video File",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Video Files",
				Pattern:     "*.mp4;*.mov;*.mkv;*.avi;*.webm",
			},
		},
	})
	return filePath, err
}

func (a *App) OpenDirectoryDialog() (string, error) {
	dirPath, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select Output Folder",
	})
	return dirPath, err
}

func (a *App) GetVideoDuration(filepath string) (float64, error) {
	return processor.GetVideoDuration(filepath)
}

func (a *App) ValidateClip(clip core.Clip, totalDuration float64) (bool, error) {
	if clip.Name == "" {
		return false, fmt.Errorf("Clip Name cannot be empty")
	}

	// 1. Check Format
	startSeconds, err := utils.ToSeconds(clip.StartTime)
	if err != nil {
		return false, fmt.Errorf("Invalid Start Time: %w", err)
	}
	endSeconds, err := utils.ToSeconds(clip.EndTime)
	if err != nil {
		return false, fmt.Errorf("Invalid End Time: %w", err)
	}

	// 2. Check Logic
	if startSeconds >= endSeconds {
		return false, fmt.Errorf("End Time must be after Start Time")
	}

	// 3. Check Range
	if endSeconds > totalDuration {
		return false, fmt.Errorf("End Time (%.2fs) is past the video's total duration (%.2fs)", endSeconds, totalDuration)
	}

	// If all checks pass:
	return true, nil
}

// StartCropping delegates to processor package
// It passes its context so the processor can emit events.
func (a *App) StartCropping(clips []core.Clip, inputPath string, outputDir string) {
	// This call is non-blocking. It returns immediately,
	// and the processor.StartCropping func runs in the background.
	processor.StartCropping(a.ctx, clips, inputPath, outputDir)
}
