package ai

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// RunWhisper invoca el script de Python que usa faster-whisper.
func RunWhisper(input, output, model, device string) error {
	py := pythonExecPath()
	script := filepath.Join("aiworker", "transcribe.py")
	args := []string{"-u", script, "--input", input, "--output", output, "--model", model, "--device", device}
	cmd := exec.Command(py, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func pythonExecPath() string {
	// 1) Permitir override por env var
	if p := os.Getenv("VENV_PYTHON"); p != "" {
		return p
	}
	// 2) Buscar venv local
	if runtime.GOOS == "windows" {
		p := filepath.Join("venv", "Scripts", "python.exe")
		if _, err := os.Stat(p); err == nil {
			return p
		}
		return "python"
	}
	p := filepath.Join("venv", "bin", "python3")
	if _, err := os.Stat(p); err == nil {
		return p
	}
	// fallback
	return "python3"
}
