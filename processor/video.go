package processor

import (
	"context"
	"cropper/core"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// GetVideoDuration uses ffprobe to read the video's duration in seconds
func GetVideoDuration(filepath string) (float64, error) {
	cmd := exec.Command("ffprobe",
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		filepath,
	)

	output, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("failed to run ffprobe: %w. Make sure ffprobe is in your PATH.", err)
	}

	durationStr := strings.TrimSpace(string(output))
	if durationStr == "" {
		return 0, fmt.Errorf("could not get duration from file: %s", filepath)
	}

	duration, err := strconv.ParseFloat(durationStr, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse duration: %w", err)
	}

	return duration, nil
}

// StartCropping begins the extraction process on a background go-routine
// It accepts the app context so it can emit events back to the UI.
func StartCropping(ctx context.Context, clips []core.Clip, inputPath string, outputDir string) {

	go func() {
		totalClips := len(clips)
		var timePerClip float64 = 0 // For ETA
		var startTime time.Time

		for i, clip := range clips {
			startTime = time.Now()

			//1. Send Status Update
			status := fmt.Sprintf("Processing clip %d of %d: %s", i+1, totalClips, clip.Name)
			runtime.EventsEmit(ctx, "processing:status", status)

			//2.  Build FFmpeg Command
			// Output file is named based on clip name
			outputFile := fmt.Sprintf("%s/%s.mp4", outputDir, clip.Name)

			// -i = input file, -ss = start time, -to = end time
			// -c copy = "stream copy".
			// -y = overwrite output file without asking
			cmd := exec.Command("ffmpeg",
				"-i", inputPath,
				"-ss", clip.StartTime,
				"-to", clip.EndTime,
				"-c", "copy",
				"-y",
				outputFile,
			)

			//3.  Run Command and Handle Errors
			output, err := cmd.CombinedOutput()
			if err != nil {
				// Send detailed error back to UI
				errorMsg := fmt.Sprintf("Error on clip '%s': %s", clip.Name, string(output))
				runtime.EventsEmit(ctx, "processing:error", errorMsg)
				return // Stop processing if one clip fails
			}

			//4. Send Progress & ETA Updates
			// Progress Bar update
			progress := float64(i+1) / float64(totalClips)
			runtime.EventsEmit(ctx, "processing:progress", progress)

			// Time Estimation (ETA) calculation
			if i == 0 {
				// First clip, set the average
				timePerClip = time.Since(startTime).Seconds()
			} else {
				// Update a rolling average for better accuracy
				timePerClip = (timePerClip + time.Since(startTime).Seconds()) / 2
			}

			clipsRemaining := totalClips - (i + 1)
			etaSeconds := timePerClip * float64(clipsRemaining)
			etaDuration := time.Duration(etaSeconds) * time.Second

			// Send ETA string
			runtime.EventsEmit(ctx, "processing:eta", etaDuration.Round(time.Second).String())
		}

		//5. Send Final "Done" Message
		runtime.EventsEmit(ctx, "processing:status", "All clips created successfully!")
		runtime.EventsEmit(ctx, "processing:eta", "Done")
	}()
}
