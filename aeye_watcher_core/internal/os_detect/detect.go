package os_detect

import (
	"os/exec"
	"runtime"
	"strings"
)

type OSInfo struct {
	OS           string `json:"os"`
	Arch         string `json:"arch"`
	Distro       string `json:"distro,omitempty"`
	Version      string `json:"version,omitempty"`
	Hostname     string `json:"hostname,omitempty"`
	IsContainer  bool   `json:"is_container"`
	GPUAvailable bool   `json:"gpu_available"`
}

func Detect() OSInfo {
	info := OSInfo{
		OS:   runtime.GOOS,
		Arch: runtime.GOARCH,
	}

	switch runtime.GOOS {
	case "linux":
		info.Distro, info.Version = detectLinuxDistro()
		info.IsContainer = isContainer()
	case "darwin":
		info.Version = detectMacOSVersion()
		info.Distro = "macOS"
	case "windows":
		info.Version = detectWindowsVersion()
		info.Distro = "Windows"
	}

	info.GPUAvailable = checkGPU()
	info.Hostname = getHostname()
	return info
}

func detectLinuxDistro() (distro, version string) {
	out, err := exec.Command("cat", "/etc/os-release").Output()
	if err != nil {
		return "Linux", "Unknown"
	}
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "NAME=") {
			distro = strings.Trim(strings.TrimPrefix(line, "NAME="), `"`)
		}
		if strings.HasPrefix(line, "VERSION_ID=") {
			version = strings.Trim(strings.TrimPrefix(line, "VERSION_ID="), `"`)
		}
	}
	if distro == "" {
		distro = "Linux"
	}
	return
}

func detectMacOSVersion() string {
	out, err := exec.Command("sw_vers", "-productVersion").Output()
	if err != nil {
		return "Unknown"
	}
	return strings.TrimSpace(string(out))
}

func detectWindowsVersion() string {
	out, err := exec.Command("cmd", "/c", "ver").Output()
	if err != nil {
		return "Unknown"
	}
	return strings.TrimSpace(string(out))
}

func isContainer() bool {
	_, err1 := exec.LookPath("docker")
	_, err2 := exec.Command("cat", "/.dockerenv").Output()
	return err1 == nil || err2 == nil
}

func checkGPU() bool {
	_, err := exec.LookPath("nvidia-smi")
	if err == nil {
		return true
	}
	_, err = exec.LookPath("rocm-smi")
	return err == nil
}

func getHostname() string {
	out, err := exec.Command("hostname").Output()
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(string(out))
}
