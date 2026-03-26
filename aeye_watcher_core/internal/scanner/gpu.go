package scanner

import (
	"bytes"
	"os/exec"
	"strconv"
	"strings"

	"github.com/aeye/watcher-core/internal/models"
)

func ScanGPU() models.GPUInfo {
	// Try NVIDIA first
	gpu := tryNvidia()
	if gpu.Error == "" {
		return gpu
	}

	// Try AMD ROCm
	gpu = tryAMD()
	if gpu.Error == "" {
		return gpu
	}

	// Try Apple Silicon via system_profiler
	gpu = tryAppleSilicon()
	if gpu.Error == "" {
		return gpu
	}

	return models.GPUInfo{Error: "No GPU or GPU monitoring tools found", Backend: "none"}
}

func tryNvidia() models.GPUInfo {
	cmd := exec.Command("nvidia-smi",
		"--query-gpu=name,memory.used,memory.total,temperature.gpu",
		"--format=csv,noheader,nounits")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return models.GPUInfo{Error: "nvidia-smi not available"}
	}

	parts := strings.Split(strings.TrimSpace(out.String()), ",")
	if len(parts) < 3 {
		return models.GPUInfo{Error: "Unexpected nvidia-smi output"}
	}

	info := models.GPUInfo{Backend: "nvidia"}
	info.Name = strings.TrimSpace(parts[0])
	info.UsedMB, _ = strconv.Atoi(strings.TrimSpace(parts[1]))
	info.TotalMB, _ = strconv.Atoi(strings.TrimSpace(parts[2]))
	if len(parts) >= 4 {
		info.TempC, _ = strconv.Atoi(strings.TrimSpace(parts[3]))
	}
	if info.TotalMB > 0 {
		info.UsePct = float64(info.UsedMB) / float64(info.TotalMB) * 100
	}
	return info
}

func tryAMD() models.GPUInfo {
	cmd := exec.Command("rocm-smi", "--showmeminfo", "vram", "--json")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return models.GPUInfo{Error: "rocm-smi not available"}
	}
	// Basic parsing - ROCm JSON varies by version
	return models.GPUInfo{Backend: "amd", Name: "AMD GPU (ROCm)", Error: ""}
}

func tryAppleSilicon() models.GPUInfo {
	// Apple Silicon shares RAM; we can get system memory usage as proxy
	cmd := exec.Command("sysctl", "-n", "hw.memsize")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return models.GPUInfo{Error: "Not Apple Silicon or sysctl unavailable"}
	}
	totalBytes, _ := strconv.ParseInt(strings.TrimSpace(out.String()), 10, 64)
	totalMB := int(totalBytes / 1024 / 1024)
	return models.GPUInfo{
		Backend: "apple",
		Name:    "Apple Silicon (Unified Memory)",
		TotalMB: totalMB,
		UsedMB:  0, // Would need additional tooling
		UsePct:  0,
	}
}
