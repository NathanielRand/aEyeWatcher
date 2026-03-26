package scanner

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/aeye/watcher-core/internal/models"
)

type AgentEcosystemDef struct {
	Name        string
	Ecosystem   string
	Description string
	Icon        string
	CLIName     string
	Expected    []models.ExpectedFile
}

func getAgentDefs(home string) []AgentEcosystemDef {
	return []AgentEcosystemDef{
		{
			Name: "Claude Code", Ecosystem: "claude", Description: "Anthropic's agentic coding assistant", Icon: "🧠", CLIName: "claude",
			Expected: []models.ExpectedFile{
				{Label: "CLI Binary", Pattern: "claude"},
				{Label: "Desktop Config", Pattern: filepath.Join(home, ".claude", "claude_desktop_config.json")},
				{Label: "Agent Settings", Pattern: filepath.Join(home, ".claude", "settings.json")},
				{Label: "Project CLAUDE.md", Pattern: "CLAUDE.md"},
				{Label: "Memory / SUMMARY", Pattern: filepath.Join(home, ".claude", "SUMMARY.md")},
				{Label: "MCP Config", Pattern: filepath.Join(home, ".claude", "mcp.json")},
			},
		},
		{
			Name: "OpenInterpreter", Ecosystem: "openinterpreter", Description: "Open-source code interpreter agent", Icon: "🐍", CLIName: "interpreter",
			Expected: []models.ExpectedFile{
				{Label: "CLI Binary", Pattern: "interpreter"},
				{Label: "Config File", Pattern: filepath.Join(home, ".openinterpreter", "config.yaml")},
				{Label: "Profile", Pattern: filepath.Join(home, ".openinterpreter", "profile.yaml")},
				{Label: "History", Pattern: filepath.Join(home, ".openinterpreter", "history.db")},
				{Label: "Custom Skills", Pattern: filepath.Join(home, ".openinterpreter", "skills")},
			},
		},
		{
			Name: "Aider", Ecosystem: "aider", Description: "AI pair programmer in your terminal", Icon: "✍️", CLIName: "aider",
			Expected: []models.ExpectedFile{
				{Label: "CLI Binary", Pattern: "aider"},
				{Label: "Config (.aider.conf)", Pattern: filepath.Join(home, ".aider.conf.yml")},
				{Label: ".aiderignore", Pattern: ".aiderignore"},
				{Label: "CONVENTIONS.md", Pattern: "CONVENTIONS.md"},
				{Label: "History File", Pattern: filepath.Join(home, ".aider.chat.history.md")},
				{Label: "Model Metadata", Pattern: filepath.Join(home, ".aider.model.metadata.json")},
			},
		},
		{
			Name: "AutoGPT", Ecosystem: "autogpt", Description: "Autonomous AI agent framework", Icon: "🤖", CLIName: "autogpt",
			Expected: []models.ExpectedFile{
				{Label: "CLI Binary", Pattern: "autogpt"},
				{Label: ".env Config", Pattern: ".env"},
				{Label: "AI Settings", Pattern: "ai_settings.yaml"},
				{Label: "Plugin Dir", Pattern: "plugins"},
				{Label: "Memory Dir", Pattern: "auto_gpt_workspace"},
				{Label: "Prompts Config", Pattern: filepath.Join("autogpt", "prompts")},
			},
		},
		{
			Name: "CrewAI", Ecosystem: "crewai", Description: "Multi-agent orchestration framework", Icon: "🚢", CLIName: "crewai",
			Expected: []models.ExpectedFile{
				{Label: "CLI Binary", Pattern: "crewai"},
				{Label: "Crew Config", Pattern: "crew.yaml"},
				{Label: "Agents Config", Pattern: filepath.Join("config", "agents.yaml")},
				{Label: "Tasks Config", Pattern: filepath.Join("config", "tasks.yaml")},
				{Label: "Tools Dir", Pattern: "tools"},
				{Label: "Knowledge Base", Pattern: "knowledge"},
			},
		},
		{
			Name: "LangChain", Ecosystem: "langchain", Description: "LLM application framework", Icon: "🔗", CLIName: "langchain",
			Expected: []models.ExpectedFile{
				{Label: "CLI Binary", Pattern: "langchain"},
				{Label: "LangSmith Config", Pattern: filepath.Join(home, ".langchain", "config.json")},
				{Label: "Chains Dir", Pattern: "chains"},
				{Label: "Tools Dir", Pattern: "tools"},
				{Label: "Vectorstore", Pattern: "vectorstore"},
				{Label: ".env (API Keys)", Pattern: ".env"},
			},
		},
		{
			Name: "Ollama", Ecosystem: "ollama", Description: "Local model runner & server", Icon: "🦙", CLIName: "ollama",
			Expected: []models.ExpectedFile{
				{Label: "CLI Binary", Pattern: "ollama"},
				{Label: "Models Directory", Pattern: filepath.Join(home, ".ollama", "models")},
				{Label: "Config File", Pattern: filepath.Join(home, ".ollama", "config")},
				{Label: "Manifests", Pattern: filepath.Join(home, ".ollama", "models", "manifests")},
				{Label: "Blobs / Weights", Pattern: filepath.Join(home, ".ollama", "models", "blobs")},
				{Label: "Modelfile", Pattern: "Modelfile"},
			},
		},
		{
			Name: "LM Studio", Ecosystem: "lmstudio", Description: "Desktop LLM UI and server", Icon: "🖥️", CLIName: "lms",
			Expected: []models.ExpectedFile{
				{Label: "CLI (lms)", Pattern: "lms"},
				{Label: "Models Cache (Linux)", Pattern: filepath.Join(home, ".cache", "lm-studio", "models")},
				{Label: "Models Cache (Linux alt)", Pattern: filepath.Join(home, ".local", "share", "lm-studio", "models")},
				{Label: "Config Dir", Pattern: filepath.Join(home, ".lmstudio")},
				{Label: "Plugin Dir", Pattern: filepath.Join(home, ".lmstudio", "plugins")},
				{Label: "Presets Dir", Pattern: filepath.Join(home, ".lmstudio", "presets")},
			},
		},
		{
			Name: "Open WebUI", Ecosystem: "openwebui", Description: "Web interface for local LLMs", Icon: "🌐", CLIName: "open-webui",
			Expected: []models.ExpectedFile{
				{Label: "CLI Binary", Pattern: "open-webui"},
				{Label: "Data Dir", Pattern: filepath.Join(home, ".open-webui")},
				{Label: "Config", Pattern: filepath.Join(home, ".open-webui", "config.json")},
				{Label: "Database", Pattern: filepath.Join(home, ".open-webui", "webui.db")},
				{Label: "Documents Dir", Pattern: filepath.Join(home, ".open-webui", "docs")},
				{Label: "Tools Dir", Pattern: filepath.Join(home, ".open-webui", "tools")},
			},
		},
		{
			Name: "Continue", Ecosystem: "continue", Description: "VS Code AI coding assistant", Icon: "▶️", CLIName: "continue",
			Expected: []models.ExpectedFile{
				{Label: "Config File", Pattern: filepath.Join(home, ".continue", "config.json")},
				{Label: "Config TS", Pattern: filepath.Join(home, ".continue", "config.ts")},
				{Label: "Prompts Dir", Pattern: filepath.Join(home, ".continue", "prompts")},
				{Label: "Rules Dir", Pattern: filepath.Join(home, ".continue", "rules")},
				{Label: "Slash Commands", Pattern: filepath.Join(home, ".continue", "config.json")},
				{Label: "Docs", Pattern: filepath.Join(home, ".continue", "docs")},
			},
		},
	}
}

func resolveExpectedFile(home, cwd, pattern, cliName string) (found bool, path string) {
	// 1. Try as absolute path
	if filepath.IsAbs(pattern) {
		if _, err := os.Stat(pattern); err == nil {
			return true, pattern
		}
		// Check dir existence
		if _, err := os.Stat(pattern); err == nil {
			return true, pattern
		}
		return false, ""
	}

	// 2. CLI binary check
	if !strings.Contains(pattern, string(os.PathSeparator)) && !strings.Contains(pattern, ".") {
		// Looks like a CLI name
		if p, err := findExecutable(pattern); err == nil {
			return true, p
		}
		// Fall through to file check
	}

	// 3. Try relative to cwd
	cwdPath := filepath.Join(cwd, pattern)
	if _, err := os.Stat(cwdPath); err == nil {
		return true, cwdPath
	}

	// 4. Try relative to home
	homePath := filepath.Join(home, pattern)
	if _, err := os.Stat(homePath); err == nil {
		return true, homePath
	}

	return false, ""
}

func findExecutable(name string) (string, error) {
	// Check PATH
	return "", os.ErrNotExist
}

func BuildAgentCards(home string, progressCh chan<- string) []models.AgentCard {
	cwd, _ := os.Getwd()
	defs := getAgentDefs(home)
	var cards []models.AgentCard

	for _, def := range defs {
		if progressCh != nil {
			progressCh <- "Profiling " + def.Name + " ecosystem..."
		}

		card := models.AgentCard{
			Name:        def.Name,
			Ecosystem:   def.Ecosystem,
			Description: def.Description,
			Icon:        def.Icon,
			Status:      "not_found",
		}

		foundCount := 0
		for _, ef := range def.Expected {
			found, path := resolveExpectedFile(home, cwd, ef.Pattern, def.CLIName)
			ef.Found = found
			ef.Path = path
			if found {
				foundCount++
			}
			card.ExpectedFiles = append(card.ExpectedFiles, ef)
		}

		if foundCount == len(def.Expected) {
			card.Status = "detected"
		} else if foundCount > 0 {
			card.Status = "partial"
		}

		// Also scan for loose agent files
		card.Files = scanAgentFiles(home, def.Ecosystem)
		cards = append(cards, card)
	}

	return cards
}

func scanAgentFiles(home, ecosystem string) []models.AgentFile {
	var files []models.AgentFile

	// Map of ecosystem -> directories to look in
	ecosystemDirs := map[string][]string{
		"claude":         {filepath.Join(home, ".claude")},
		"openinterpreter": {filepath.Join(home, ".openinterpreter")},
		"aider":          {home},
		"crewai":         {".", "crewai"},
		"langchain":      {filepath.Join(home, ".langchain"), "."},
		"ollama":         {filepath.Join(home, ".ollama")},
		"continue":       {filepath.Join(home, ".continue")},
		"lmstudio":       {filepath.Join(home, ".lmstudio"), filepath.Join(home, ".cache", "lm-studio")},
	}

	dirs, ok := ecosystemDirs[ecosystem]
	if !ok {
		return files
	}

	fileTypeMap := map[string]string{
		".yaml":  "config",
		".yml":   "config",
		".json":  "config",
		".md":    "memory",
		".toml":  "config",
		".env":   "config",
		".db":    "memory",
		".gguf":  "model",
		".soul":  "soul",
		".skill": "skill",
		".tool":  "tool",
	}

	for _, dir := range dirs {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			ext := strings.ToLower(filepath.Ext(entry.Name()))
			ftype, ok := fileTypeMap[ext]
			if !ok {
				continue
			}
			info, _ := entry.Info()
			var sz int64
			if info != nil {
				sz = info.Size()
			}
			files = append(files, models.AgentFile{
				Type: ftype, Name: entry.Name(),
				Path: filepath.Join(dir, entry.Name()), SizeBytes: sz,
			})
		}
	}

	return files
}
