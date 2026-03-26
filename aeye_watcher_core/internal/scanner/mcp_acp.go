package scanner

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/aeye/watcher-core/internal/models"
)

// Known MCP config file locations per OS/agent
func getMCPConfigPaths(home string) []string {
	paths := []string{
		// Claude Desktop
		filepath.Join(home, ".claude", "claude_desktop_config.json"),
		filepath.Join(home, "Library", "Application Support", "Claude", "claude_desktop_config.json"),   // macOS
		filepath.Join(home, "AppData", "Roaming", "Claude", "claude_desktop_config.json"),                // Windows
		// Cursor
		filepath.Join(home, ".cursor", "mcp.json"),
		filepath.Join(home, ".cursor", "settings.json"),
		// VS Code / Continue
		filepath.Join(home, ".continue", "config.json"),
		// Generic MCP configs
		filepath.Join(home, ".mcp", "config.json"),
		filepath.Join(home, ".config", "mcp", "servers.json"),
		// Zed editor
		filepath.Join(home, ".config", "zed", "settings.json"),
		// Windsurf
		filepath.Join(home, ".codeium", "windsurf", "mcp_config.json"),
		// Project-level (cwd)
		".mcp.json",
		"mcp.json",
		filepath.Join(".claude", "settings.json"),
	}
	return paths
}

type claudeDesktopConfig struct {
	MCPServers map[string]struct {
		Command string            `json:"command"`
		Args    []string          `json:"args"`
		Env     map[string]string `json:"env"`
		URL     string            `json:"url"`
	} `json:"mcpServers"`
}

func ScanMCPServers(home string, progressCh chan<- string) []models.MCPServer {
	var servers []models.MCPServer
	seen := map[string]bool{}

	configPaths := getMCPConfigPaths(home)
	// Also search for any mcp*.json in home
	filepath.WalkDir(home, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			name := d.Name()
			if skipDirs[name] || name == ".ollama" {
				return filepath.SkipDir
			}
			rel, _ := filepath.Rel(home, path)
			if len(strings.Split(rel, string(os.PathSeparator))) > 5 {
				return filepath.SkipDir
			}
		}
		if !d.IsDir() {
			name := strings.ToLower(d.Name())
			if (strings.Contains(name, "mcp") || strings.Contains(name, "claude_desktop")) && strings.HasSuffix(name, ".json") {
				configPaths = append(configPaths, path)
			}
		}
		return nil
	})

	for _, cfgPath := range configPaths {
		if seen[cfgPath] {
			continue
		}
		seen[cfgPath] = true

		data, err := os.ReadFile(cfgPath)
		if err != nil {
			continue
		}

		if progressCh != nil {
			progressCh <- fmt.Sprintf("Parsing MCP config: %s", filepath.Base(cfgPath))
		}

		var cfg claudeDesktopConfig
		if err := json.Unmarshal(data, &cfg); err != nil {
			continue
		}

		for name, srv := range cfg.MCPServers {
			transport := "stdio"
			url := srv.URL
			if url != "" {
				transport = "http"
			}
			status := probeSSEEndpoint(url)
			if status == "" {
				status = "configured"
			}

			servers = append(servers, models.MCPServer{
				Name:      name,
				URL:       url,
				Config:    cfgPath,
				Status:    status,
				Transport: transport,
			})
		}
	}

	// Probe common MCP server ports
	mcpPorts := []struct {
		name string
		port string
		path string
	}{
		{"Custom MCP (3100)", "3100", "/"},
		{"Custom MCP (4000)", "4000", "/"},
		{"Custom MCP (8081)", "8081", "/"},
		{"Custom MCP (9000)", "9000", "/"},
	}

	client := http.Client{Timeout: 500 * time.Millisecond}
	for _, p := range mcpPorts {
		url := fmt.Sprintf("http://localhost:%s%s", p.port, p.path)
		resp, err := client.Get(url)
		if err == nil && resp.StatusCode < 500 {
			resp.Body.Close()
			// Check if it looks like an MCP server (SSE endpoint)
			servers = append(servers, models.MCPServer{
				Name:      p.name,
				URL:       url,
				Status:    "running",
				Transport: "http",
			})
		}
	}

	// Windows-specific paths
	if runtime.GOOS == "windows" {
		winPaths := []string{
			filepath.Join(home, "AppData", "Roaming", "Code", "User", "settings.json"),
		}
		for _, p := range winPaths {
			if _, err := os.Stat(p); err == nil {
				configPaths = append(configPaths, p)
			}
		}
	}

	return servers
}

func probeSSEEndpoint(url string) string {
	if url == "" {
		return ""
	}
	client := http.Client{Timeout: 800 * time.Millisecond}
	resp, err := client.Get(url)
	if err != nil {
		return "unreachable"
	}
	defer resp.Body.Close()
	if resp.StatusCode < 400 {
		return "running"
	}
	return "error"
}

// ACP (Agent Communication Protocol) scanner
var acpPorts = []struct {
	port string
	name string
}{
	{"8765", "ACP Agent (8765)"},
	{"9999", "ACP Agent (9999)"},
	{"7000", "ACP Agent (7000)"},
	{"5000", "ACP Agent (5000)"},
}

func ScanACPAgents(progressCh chan<- string) []models.ACPAgent {
	var agents []models.ACPAgent
	client := http.Client{Timeout: 600 * time.Millisecond}

	for _, p := range acpPorts {
		// ACP uses /.well-known/agent.json as discovery endpoint (per spec)
		url := fmt.Sprintf("http://localhost:%s/.well-known/agent.json", p.port)
		resp, err := client.Get(url)
		if err != nil {
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == 200 {
			var agentMeta struct {
				Name string `json:"name"`
			}
			json.NewDecoder(resp.Body).Decode(&agentMeta)
			name := agentMeta.Name
			if name == "" {
				name = p.name
			}
			agents = append(agents, models.ACPAgent{
				Name:     name,
				Endpoint: fmt.Sprintf("http://localhost:%s", p.port),
				Status:   "running",
				Protocol: "ACP",
			})
		}
	}
	return agents
}
