package video

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type VideoMediaEncoder struct {
	MaxClipDuration int
	ffmpegPath      string
	ffmprobePath    string
}

func NewVideoMediaEncoder(ffmpegPath string, ffmprobePath string, maxDuration int) *VideoMediaEncoder {
	return &VideoMediaEncoder{
		ffmpegPath:      ffmpegPath,
		MaxClipDuration: maxDuration,
	}
}

func (v *VideoMediaEncoder) GenerateClips(ctx context.Context, videoPath string, outputDir string) (string, error) {
	err := os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		return "", err
	}

	// Get video duration using ffprobe
	durationCommands := []string{"-v", "error", "-show_entries", "format=duration", "-of", "default=noprint_wrappers=1:nokey=1", videoPath}
	log.Printf("Running command: %s %s\n", "ffprobe", strings.Join(durationCommands, " "))
	durationCmd := exec.CommandContext(ctx, "ffprobe", durationCommands...)
	durationOutput, err := durationCmd.Output()
	if err != nil {
		return "", errors.Wrap(err, "failed to get video duration for clips generation")
	}

	durationStr := strings.TrimSpace(string(durationOutput))
	videoDuration, err := strconv.ParseFloat(durationStr, 64)
	if err != nil {
		return "", fmt.Errorf("failed to parse duration as float: %v", err)
	}

	// Calculate the number of clips based on the video duration
	numClips := int((videoDuration - 35) / 20)
	if numClips > 9 {
		numClips = 9
	}
	if numClips < 1 {
		return "", fmt.Errorf("video is too short to generate any clips")
	}

	rand.Seed(time.Now().UnixNano())
	var startTime float64

	for i := 0; i < numClips; i++ {
		// Generate random start time within the first 20 seconds
		if i == 0 {
			startTime = rand.Float64() * 20
		} else {
			// Generate random start time with a dynamic maximum
			maxRandomStart := (videoDuration * 2 / float64(numClips)) - (float64(numClips-i) + 5)
			startTime += 10 + rand.Float64()*(maxRandomStart-10)
		}
		if startTime < 0 {
			startTime = 0
		}
		if videoDuration-10 < startTime {
			break
		}
		log.Println(startTime)
		// Generate random duration for the clip (up to 10 seconds)
		clipDuration := rand.Float64()*10 + float64(v.MaxClipDuration)

		// Generate output file path
		outputFile := filepath.Join(outputDir, fmt.Sprintf("clip%d.mp4", time.Now().UnixNano()))

		cmd := exec.Command(
			v.ffmpegPath,
			"-ss", fmt.Sprintf("%.2f", startTime),
			"-i", videoPath,
			"-t", fmt.Sprintf("%.2f", clipDuration),
			"-vf", "crop=ih*9/16:ih", // Scale to height 418 and then crop to 19.5:9 aspect ratio
			"-c:a", "copy",
			outputFile,
		)
		// )
		if err := cmd.Run(); err != nil {
			return "", fmt.Errorf("failed to run FFMPEG: %v", err)
		}

		log.Printf("Random clip generated: %s", outputFile)
	}

	return outputDir, nil
}
