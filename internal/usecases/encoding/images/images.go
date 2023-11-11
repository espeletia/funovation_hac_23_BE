package images

import (
	"bytes"
	"context"
	"fmt"
	"funovation_23/internal/domain"
	"os/exec"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type ImageMediaEncoder struct {
	ffmpegPath string
}

func NewImageMediaEncoder(ffmpegPath string) *ImageMediaEncoder {
	return &ImageMediaEncoder{ffmpegPath: ffmpegPath}
}

// ffmpeg -i "videoInput" -ss "$(shuf -i 1-$(printf "%.0f" $(ffprobe -v error -show_entries format=duration -of default=noprint_wrappers=1:nokey=1 "videoInput")) -n 1)" -vframes 1 "oitputJPG"

func (v *ImageMediaEncoder) GenerateThumbanail(ctx context.Context, srcPath, tempDir string) error {
	durationCmd := exec.CommandContext(ctx, "ffprobe", "-v", "error", "-show_entries", "format=duration", "-of", "default=noprint_wrappers=1:nokey=1", srcPath)
	durationOutput, err := durationCmd.Output()
	if err != nil {
		return errors.Wrap(err, "failed to get video duration for thumbnail generation")
	}

	durationStr := strings.TrimSpace(string(durationOutput))

	durationFloat, err := strconv.ParseFloat(strings.TrimSpace(durationStr), 64)
	if err != nil {
		return errors.Wrap(err, "failed to parse duration as float")
	}

	durationInt := int64(durationFloat)
	if durationInt <= 0 {
		return domain.UnableToThumbnail
	}
	// Generate a random timestamp using shuf
	randomTimeCmd := exec.CommandContext(ctx, "shuf", "-i", fmt.Sprintf("1-%s", strconv.FormatInt(durationInt, 10)), "-n", "1")
	randomTimeOutput, err := randomTimeCmd.Output()
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to generate random timestamp: %s", randomTimeCmd.Stderr))
	}

	randomTimeStr := strings.TrimSpace(string(randomTimeOutput))

	// Build the final ffmpeg command
	command := []string{
		"-i",
		srcPath,
		"-ss",
		randomTimeStr,
		"-vframes",
		"1",
		tempDir,
	}

	// log.Println(command)

	// Execute the ffmpeg command
	cmd := exec.CommandContext(ctx, "ffmpeg", command...)
	errorBuffer := bytes.Buffer{}
	cmd.Stderr = &errorBuffer
	cmd.Stdout = &errorBuffer
	err = cmd.Run()
	if err != nil {
		return errors.Wrap(err, errorBuffer.String())
	}
	return nil
}
