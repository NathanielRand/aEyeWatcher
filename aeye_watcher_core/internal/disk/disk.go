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
	out, err := exec.Command("df", "-BGB", "--output=source,fstype,size,used,avail,pcent,target").Output()
	if err != nil {
		// fallback: basic df
		out, err = exec.Command("df", "-h").Output()
		if err != nil {
			return disks
		}
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	for i, line := range lines {
		if i == 0 {
			continue // skip header
		}
		fields := strings.Fields(line)
		if len(fields) < 6 {
			continue
		}

		device := fields[0]
		// Skip pseudo filesystems
		if strings.HasPrefix(device, "tmpfs") || strings.HasPrefix(device, "devtmpfs") ||
			strings.HasPrefix(device, "udev") || strings.HasPrefix(device, "cgroupfs") ||
			strings.HasPrefix(device, "overlay") || device == "none" {
			continue
		}

		di := models.DiskInfo{}
		di.Device = device
		if len(fields) >= 7 {
			di.FSType = fields[1]
			di.TotalGB = parseGBField(fields[2])
			di.UsedGB = parseGBField(fields[3])
			di.FreeGB = parseGBField(fields[4])
			di.UsePct = parsePct(fields[5])
			di.Mount = fields[6]
		} else {
			// basic df -h fallback
			di.TotalGB = parseHumanSize(fields[1])
			di.UsedGB = parseHumanSize(fields[2])
			di.FreeGB = parseHumanSize(fields[3])
			di.UsePct = parsePct(fields[4])
			di.Mount = fields[5]
		}

		// Only include real block devices or important mounts
		if strings.HasPrefix(device, "/dev/") || di.Mount == "/" || strings.HasPrefix(di.Mount, "/home") || strings.HasPrefix(di.Mount, "/mnt") {
			disks = append(disks, di)
		}
	}
	return disks
}

func scanWindows() []models.DiskInfo {
	var disks []models.DiskInfo
	out, err := exec.Command("wmic", "logicaldisk", "get", "DeviceID,Size,FreeSpace,FileSystem", "/format:csv").Output()
	if err != nil {
		return disks
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "Node") {
			continue
		}
		parts := strings.Split(line, ",")
		if len(parts) < 5 {
			continue
		}
		totalBytes, _ := strconv.ParseInt(strings.TrimSpace(parts[3]), 10, 64)
		freeBytes, _ := strconv.ParseInt(strings.TrimSpace(parts[2]), 10, 64)
		if totalBytes == 0 {
			continue
		}
		totalGB := float64(totalBytes) / 1e9
		freeGB := float64(freeBytes) / 1e9
		usedGB := totalGB - freeGB
		disks = append(disks, models.DiskInfo{
			Device:  strings.TrimSpace(parts[1]),
			Mount:   strings.TrimSpace(parts[1]),
			FSType:  strings.TrimSpace(parts[4]),
			TotalGB: totalGB,
			UsedGB:  usedGB,
			FreeGB:  freeGB,
			UsePct:  usedGB / totalGB * 100,
		})
	}
	return disks
}

func parseGBField(s string) float64 {
	s = strings.TrimSuffix(s, "GB")
	s = strings.TrimSuffix(s, "G")
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

func parseHumanSize(s string) float64 {
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		return 0
	}
	suffix := string(s[len(s)-1])
	num := s[:len(s)-1]
	f, err := strconv.ParseFloat(num, 64)
	if err != nil {
		return 0
	}
	switch strings.ToUpper(suffix) {
	case "T":
		return f * 1024
	case "G":
		return f
	case "M":
		return f / 1024
	case "K":
		return f / (1024 * 1024)
	}
	// It might be pure bytes if no suffix matched
	n, _ := strconv.ParseFloat(s, 64)
	return n / 1e9
}

func parsePct(s string) float64 {
	s = strings.TrimSuffix(s, "%")
	f, _ := strconv.ParseFloat(strings.TrimSpace(s), 64)
	return f
}

// FormatGB returns a human-readable string for disk sizes
func FormatGB(gb float64) string {
	if gb >= 1024 {
		return fmt.Sprintf("%.1f TB", gb/1024)
	}
	return fmt.Sprintf("%.1f GB", gb)
}
