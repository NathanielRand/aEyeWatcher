package disk

import (
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/aeye/watcher-core/internal/models"
)

func ScanDisks() []models.DiskInfo {
	switch runtime.GOOS {
	case "linux", "darwin":
		return scanUnix()
	case "windows":
		return scanWindows()
	default:
		return nil
	}
}

func scanUnix() []models.DiskInfo {
	var disks []models.DiskInfo
	// Use plain bytes (-k gives KB, easier to parse without locale issues)
	out, err := exec.Command("df", "-k").Output()
	if err != nil {
		return disks
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	for i, line := range lines {
		if i == 0 {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 6 {
			continue
		}
		device := fields[0]
		if strings.HasPrefix(device, "tmpfs") || strings.HasPrefix(device, "devtmpfs") ||
			strings.HasPrefix(device, "udev") || strings.HasPrefix(device, "cgroupfs") ||
			strings.HasPrefix(device, "overlay") || device == "none" ||
			strings.HasPrefix(device, "rootfs") {
			continue
		}
		// df -k: Filesystem 1K-blocks Used Available Use% Mounted
		totalKB, _ := strconv.ParseInt(fields[1], 10, 64)
		usedKB, _ := strconv.ParseInt(fields[2], 10, 64)
		availKB, _ := strconv.ParseInt(fields[3], 10, 64)
		pct := parsePct(fields[4])
		mount := fields[5]

		if totalKB == 0 {
			continue
		}
		if !strings.HasPrefix(device, "/dev/") && mount != "/" &&
			!strings.HasPrefix(mount, "/home") && !strings.HasPrefix(mount, "/mnt") &&
			!strings.HasPrefix(mount, "/Volumes") {
			continue
		}
		disks = append(disks, models.DiskInfo{
			Device:  device,
			Mount:   mount,
			TotalGB: float64(totalKB) / 1024 / 1024,
			UsedGB:  float64(usedKB) / 1024 / 1024,
			FreeGB:  float64(availKB) / 1024 / 1024,
			UsePct:  pct,
		})
	}
	return disks
}

func scanWindows() []models.DiskInfo {
	var disks []models.DiskInfo
	// Use PowerShell Get-PSDrive for accurate byte counts — wmic is deprecated on Win11
	out, err := exec.Command("powershell", "-NoProfile", "-Command",
		"Get-PSDrive -PSProvider FileSystem | Select-Object Name,Used,Free | ConvertTo-Csv -NoTypeInformation").Output()
	if err == nil && len(out) > 10 {
		return parsePSDrive(string(out))
	}
	// Fallback: wmic (Win10)
	out, err = exec.Command("wmic", "logicaldisk", "get",
		"DeviceID,Size,FreeSpace,FileSystem", "/format:csv").Output()
	if err != nil {
		return disks
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(strings.ReplaceAll(line, "\r", ""))
		if line == "" || strings.HasPrefix(line, "Node") || strings.HasPrefix(line, "DeviceID") {
			continue
		}
		parts := strings.Split(line, ",")
		// CSV from wmic: Node,DeviceID,FileSystem,FreeSpace,Size
		if len(parts) < 5 {
			continue
		}
		deviceID := strings.TrimSpace(parts[1])
		fsType := strings.TrimSpace(parts[2])
		freeBytes, _ := strconv.ParseInt(strings.TrimSpace(parts[3]), 10, 64)
		totalBytes, _ := strconv.ParseInt(strings.TrimSpace(parts[4]), 10, 64)
		if totalBytes == 0 {
			continue
		}
		totalGB := float64(totalBytes) / 1e9
		freeGB := float64(freeBytes) / 1e9
		usedGB := totalGB - freeGB
		disks = append(disks, models.DiskInfo{
			Device:  deviceID,
			Mount:   deviceID,
			FSType:  fsType,
			TotalGB: totalGB,
			UsedGB:  usedGB,
			FreeGB:  freeGB,
			UsePct:  usedGB / totalGB * 100,
		})
	}
	return disks
}

func parsePSDrive(csv string) []models.DiskInfo {
	var disks []models.DiskInfo
	lines := strings.Split(strings.TrimSpace(csv), "\n")
	for i, line := range lines {
		if i == 0 {
			continue // header
		}
		line = strings.TrimSpace(strings.ReplaceAll(line, "\r", ""))
		// Strip quotes
		line = strings.ReplaceAll(line, `"`, "")
		parts := strings.Split(line, ",")
		if len(parts) < 3 {
			continue
		}
		name := strings.TrimSpace(parts[0])
		usedBytes, _ := strconv.ParseInt(strings.TrimSpace(parts[1]), 10, 64)
		freeBytes, _ := strconv.ParseInt(strings.TrimSpace(parts[2]), 10, 64)
		totalBytes := usedBytes + freeBytes
		if totalBytes == 0 {
			continue
		}
		totalGB := float64(totalBytes) / 1e9
		usedGB := float64(usedBytes) / 1e9
		freeGB := float64(freeBytes) / 1e9
		disks = append(disks, models.DiskInfo{
			Device:  name + ":",
			Mount:   name + ":",
			TotalGB: totalGB,
			UsedGB:  usedGB,
			FreeGB:  freeGB,
			UsePct:  usedGB / totalGB * 100,
		})
	}
	return disks
}

func parsePct(s string) float64 {
	s = strings.TrimSuffix(strings.TrimSpace(s), "%")
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

func FormatGB(gb float64) string {
	if gb >= 1024 {
		return fmt.Sprintf("%.1f TB", gb/1024)
	}
	return fmt.Sprintf("%.1f GB", gb)
}
