package ffmpeg

import (
	"runtime"
)

// SimpleExportArgs crea una línea de comandos mínima con burn-in de ASS si se provee.
func SimpleExportArgs(input, ass, output string) []string {
	args := []string{"-y", "-hide_banner", "-loglevel", "error", "-i", input}
	filter := []string{}
	if ass != "" {
		filter = append(filter, "subtitles="+ass)
	}
	if len(filter) > 0 {
		args = append(args, "-vf", joinFilters(filter))
	}

	// Codec según plataforma (placeholder sensato)
	switch runtime.GOOS {
	case "windows":
		args = append(args, "-c:v", "h264_nvenc", "-cq", "23", "-preset", "p5")
	case "darwin":
		args = append(args, "-c:v", "h264_videotoolbox")
	default: // linux
		args = append(args, "-c:v", "h264_nvenc", "-cq", "23", "-preset", "p5")
	}
	args = append(args, "-c:a", "aac", "-b:a", "192k", output)
	return args
}

func joinFilters(filters []string) string {
	if len(filters) == 0 {
		return ""
	}
	out := filters[0]
	for i := 1; i < len(filters); i++ {
		out += "," + filters[i]
	}
	return out
}
