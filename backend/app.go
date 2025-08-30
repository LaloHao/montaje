package backend

import (
	"bytes"
	"context"
	"montaje/backend/ai"
	"montaje/backend/ffmpeg"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
}

// FFmpegVersion returns the first line of `ffmpeg -version`.
func (a *App) FFmpegVersion() (string, error) {
	cmd := exec.Command("ffmpeg", "-version")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return "", err
	}
	line := out.String()
	if i := bytes.IndexByte([]byte(line), '\n'); i > 0 {
		line = line[:i]
	}
	return line, nil
}

// Transcribe executes the transcribe python worker and returns the SRT path.
func (a *App) Transcribe(inputPath string) (string, error) {
	abs, err := filepath.Abs(inputPath)
	if err != nil {
		return "", err
	}
	if _, err := os.Stat(abs); err != nil {
		return "", err
	}

	outDir := filepath.Join(os.TempDir(), "montaje")
	_ = os.MkdirAll(outDir, 0o755)
	outSRT := filepath.Join(outDir, filepath.Base(abs)+".srt")

	// El worker decide el device (cuda/metal/cpu) automáticamente
	if err := ai.RunWhisper(abs, outSRT, "large-v3", "auto"); err != nil {
		return "", err
	}
	return outSRT, nil
}

// ExportSimple placeholder for export using NVENC/VTB/QSV according to platform
func (a *App) ExportSimple(input, subtitlesASS, output string) error {
	args := ffmpeg.SimpleExportArgs(input, subtitlesASS, output)
	cmd := exec.Command("ffmpeg", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Platform returns the current OS platform
func (a *App) Platform() string { return runtime.GOOS }

// Now returns timestamp (useful for timestamps in UI)
func (a *App) Now() int64 { return time.Now().UnixMilli() }
