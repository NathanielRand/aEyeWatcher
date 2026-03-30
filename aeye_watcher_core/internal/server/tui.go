package server

import (
	"fmt"
	"strings"
	"time"

	"github.com/aeye/watcher-core/internal/models"
	"github.com/aeye/watcher-core/internal/os_detect"
)

const (
	cReset  = "\033[0m"
	cBold   = "\033[1m"
	cDim    = "\033[2m"
	cAccent = "\033[35m"  // magenta
	cGreen  = "\033[32m"
	cYellow = "\033[33m"
	cRed    = "\033[31m"
	cCyan   = "\033[36m"
	cGray   = "\033[90m"
	cWhite  = "\033[97m"
)

func PrintBanner(osInfo os_detect.OSInfo, port string) {
	fmt.Print(cAccent)
	fmt.Println(`
  ░█▀█░█▀▀░█░█░█▀▀   █░░░█░█░█▀█░▀█▀░█▀▀░█░█░█▀▀░█▀▄
  ░█▀█░█▀▀░░█░░█▀▀   █▄▄░█▀█░█░█░░█░░█░░░█▀█░█▀▀░█▀▄
  ░▀░▀░▀▀▀░░▀░░▀▀▀   ▀▀▀░▀░▀░▀░▀░░▀░░▀▀▀░▀░▀░▀▀▀░▀░▀  v2`)
	fmt.Print(cReset)
	fmt.Printf("  %s\"Every model leaves a trace. Every agent leaves a trail.\"%s\n\n", cDim, cReset)

	printBox([]string{
		fmt.Sprintf("  ◉  Core API    %s http://localhost:%s %s", cGreen, port, cReset),
		fmt.Sprintf("  ◉  OS          %s %s %s (%s) %s", cCyan, osInfo.Distro, osInfo.Version, osInfo.Arch, cReset),
		fmt.Sprintf("  ◉  Host        %s %s %s", cCyan, osInfo.Hostname, cReset),
		fmt.Sprintf("  ◉  GPU         %s %v %s", cCyan, boolStr(osInfo.GPUAvailable, "detected", "not detected"), cReset),
		fmt.Sprintf("  ◉  Container   %s %v %s", cCyan, boolStr(osInfo.IsContainer, "yes", "no"), cReset),
	})
	fmt.Println()
}

func PrintScanEvent(evt models.ScanEvent, elapsed time.Duration) {
	icon := phaseIcon(evt.Phase)
	ts := fmt.Sprintf("%s[%5.1fs]%s", cGray, elapsed.Seconds(), cReset)

	switch evt.Type {
	case models.EventCheckpoint:
		bar := progressBar(evt.Percent, 20)
		fmt.Printf("  %s  %s %s%s%s  %s\n", icon, ts, cAccent, bar, cReset, evt.Message)
	case models.EventFound:
		fmt.Printf("  %s◆%s  %s  %s%s%s\n", cGreen, cReset, ts, cGreen, evt.Message, cReset)
	case models.EventProgress:
		fmt.Printf("  %s·%s  %s  %s%s%s\n", cGray, cReset, ts, cGray, truncate(evt.Message, 70), cReset)
	case models.EventComplete:
		fmt.Println()
		PrintScanSummary(evt)
	case models.EventError:
		fmt.Printf("  %s✗%s  %s  %s%s%s\n", cRed, cReset, ts, cRed, evt.Message, cReset)
	}
}

func PrintScanSummary(evt models.ScanEvent) {
	if evt.Detail == nil {
		return
	}
	result, ok := evt.Detail.(*models.ScanResult)
	if !ok {
		return
	}

	running := 0
	for _, c := range result.Components {
		if c.Status == models.StatusRunning {
			running++
		}
	}
	orphaned := 0
	for _, m := range result.Models {
		if m.IsOrphaned {
			orphaned++
		}
	}
	var totalGB float64
	for _, m := range result.Models {
		totalGB += float64(m.SizeBytes) / 1e9
	}

	fmt.Printf("  %s┌─ SCAN COMPLETE ─────────────────────────────────────────────┐%s\n", cAccent, cReset)
	printSummaryRow("Models found",    fmt.Sprintf("%d  (%d orphaned)", len(result.Models), orphaned))
	printSummaryRow("Model storage",   fmt.Sprintf("%.1f GB", totalGB))
	printSummaryRow("Live services",   fmt.Sprintf("%d running", running))
	printSummaryRow("MCP servers",     fmt.Sprintf("%d configured", len(result.MCPServers)))
	printSummaryRow("ACP agents",      fmt.Sprintf("%d detected", len(result.ACPAgents)))
	printSummaryRow("Agent ecosystems",fmt.Sprintf("%d profiled", len(result.AgentCards)))
	printSummaryRow("Disk volumes",    fmt.Sprintf("%d mapped", len(result.Disks)))
	printSummaryRow("Scan events",     fmt.Sprintf("%d total", 0))
	fmt.Printf("  %s└──────────────────────────────────────────────────────────────┘%s\n\n", cAccent, cReset)

	// Live services
	if running > 0 {
		fmt.Printf("  %s◉ LIVE SERVICES%s\n", cGreen, cReset)
		for _, c := range result.Components {
			if c.Status == models.StatusRunning {
				portStr := ""
				if c.Port != "" {
					portStr = fmt.Sprintf(" %s:%s%s", cGray, c.Port, cReset)
				}
				fmt.Printf("    %s●%s %s%s%s%s\n", cGreen, cReset, cWhite, c.Name, cReset, portStr)
			}
		}
		fmt.Println()
	}

	// Orphaned models
	if orphaned > 0 {
		fmt.Printf("  %s⚠ ORPHANED MODELS (%d)%s\n", cYellow, orphaned, cReset)
		count := 0
		for _, m := range result.Models {
			if m.IsOrphaned && count < 5 {
				gb := float64(m.SizeBytes) / 1e9
				fmt.Printf("    %s·%s %s%s%s  %s%.1f GB%s\n", cGray, cReset, cYellow, truncate(m.Name, 45), cReset, cGray, gb, cReset)
				count++
			}
		}
		if orphaned > 5 {
			fmt.Printf("    %s... and %d more%s\n", cGray, orphaned-5, cReset)
		}
		fmt.Println()
	}
}

// ── helpers ─────────────────────────────────────────────────────────────────

func progressBar(pct, width int) string {
	filled := pct * width / 100
	if filled > width { filled = width }
	bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
	return fmt.Sprintf("[%s] %3d%%", bar, pct)
}

func phaseIcon(phase string) string {
	icons := map[string]string{
		"system": "🖥", "gpu": "⚡", "disk": "💾",
		"components": "⚙", "mcp": "⬡", "acp": "◈",
		"models": "◆", "agents": "🕵", "complete": "✓",
	}
	if icon, ok := icons[phase]; ok {
		return icon
	}
	return "·"
}

func printBox(lines []string) {
	width := 64
	fmt.Printf("  %s┌%s┐%s\n", cAccent, strings.Repeat("─", width), cReset)
	for _, line := range lines {
		fmt.Printf("  %s│%s %-*s %s│%s\n", cAccent, cReset, width-2, line, cAccent, cReset)
	}
	fmt.Printf("  %s└%s┘%s\n", cAccent, strings.Repeat("─", width), cReset)
}

func printSummaryRow(label, value string) {
	fmt.Printf("  %s│%s  %-22s %s%-35s%s  %s│%s\n",
		cAccent, cReset, label, cGreen, value, cReset, cAccent, cReset)
}

func boolStr(v bool, yes, no string) string {
	if v { return yes }
	return no
}

func truncate(s string, n int) string {
	if len(s) <= n { return s }
	return s[:n-3] + "..."
}
