package scanner

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/aeye/watcher-core/internal/models"
)

// Memory/soul file names that any agent might generate — scanned globally
var memoryFileNames = []string{
	"SOUL.md", "AGENTS.md", "USER.md", "MEMORY.md", "CONTEXT.md",
	"PERSONA.md", "IDENTITY.md", "PREFERENCES.md", "PROFILE.md",
	"CLAUDE.md", "SYSTEM.md", "INSTRUCTIONS.md", "RULES.md",
	"KNOWLEDGE.md", "SUMMARY.md", "NOTES.md", "HISTORY.md",
	".soul", ".persona", ".memory",
}

type AgentEcosystemDef struct {
	Name        string
	Ecosystem   string
	Description string
	Icon        string
	CLIName     string
	Expected    []models.ExpectedFile
}

func getAgentDefs(home string) []AgentEcosystemDef {
	// LM Studio hub path varies by OS
	lmStudioModels := lmStudioModelPath(home)

	return []AgentEcosystemDef{
		{
			Name: "Claude Code", Ecosystem: "claude", Description: "Anthropic's agentic coding assistant", Icon: "🧠", CLIName: "claude",
			Expected: []models.ExpectedFile{
				{Label: "CLI Binary", Pattern: "claude"},
				{Label: "Desktop Config", Pattern: filepath.Join(home, ".claude", "claude_desktop_config.json")},
				{Label: "Agent Settings", Pattern: filepath.Join(home, ".claude", "settings.json")},
				{Label: "CLAUDE.md (project)", Pattern: "CLAUDE.md"},
				{Label: "SOUL.md", Pattern: filepath.Join(home, ".claude", "SOUL.md")},
				{Label: "MEMORY.md", Pattern: filepath.Join(home, ".claude", "MEMORY.md")},
				{Label: "SUMMARY.md", Pattern: filepath.Join(home, ".claude", "SUMMARY.md")},
				{Label: "MCP Config", Pattern: filepath.Join(home, ".claude", "mcp.json")},
				{Label: "Skills Dir", Pattern: filepath.Join(home, ".claude", "skills")},
			},
		},
		{
			Name: "OpenClaw", Ecosystem: "openclaw", Description: "OpenClaw multi-agent framework with built-in gateway", Icon: "🦅", CLIName: "openclaw",
			Expected: []models.ExpectedFile{
				{Label: "CLI Binary", Pattern: "openclaw"},
				{Label: "Config Dir", Pattern: filepath.Join(home, ".openclaw")},
				{Label: "Config File", Pattern: filepath.Join(home, ".openclaw", "config.yaml")},
				{Label: "SOUL.md", Pattern: filepath.Join(home, ".openclaw", "SOUL.md")},
				{Label: "AGENTS.md", Pattern: filepath.Join(home, ".openclaw", "AGENTS.md")},
				{Label: "USER.md", Pattern: filepath.Join(home, ".openclaw", "USER.md")},
				{Label: "Skills Dir", Pattern: filepath.Join(home, ".openclaw", "skills")},
				{Label: "Gateway Config", Pattern: filepath.Join(home, ".openclaw", "gateway.yaml")},
				{Label: "Memory Dir", Pattern: filepath.Join(home, ".openclaw", "memory")},
			},
		},
		{
			Name: "Hermes", Ecosystem: "hermes", Description: "Hermes agent runtime with HTTP gateway", Icon: "⚡", CLIName: "hermes",
			Expected: []models.ExpectedFile{
				{Label: "CLI Binary", Pattern: "hermes"},
				{Label: "Config Dir", Pattern: filepath.Join(home, ".hermes")},
				{Label: "Config File", Pattern: filepath.Join(home, ".hermes", "config.yaml")},
				{Label: "SOUL.md", Pattern: filepath.Join(home, ".hermes", "SOUL.md")},
				{Label: "AGENTS.md", Pattern: filepath.Join(home, ".hermes", "AGENTS.md")},
				{Label: "USER.md", Pattern: filepath.Join(home, ".hermes", "USER.md")},
				{Label: "Skills Dir", Pattern: filepath.Join(home, ".hermes", "skills")},
				{Label: "Tools Dir", Pattern: filepath.Join(home, ".hermes", "tools")},
				{Label: "Personas Dir", Pattern: filepath.Join(home, ".hermes", "personas")},
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
				{Label: "USER.md", Pattern: filepath.Join(home, ".openinterpreter", "USER.md")},
			},
		},
		{
			Name: "Aider", Ecosystem: "aider", Description: "AI pair programmer in your terminal", Icon: "✍️", CLIName: "aider",
			Expected: []models.ExpectedFile{
				{Label: "CLI Binary", Pattern: "aider"},
				{Label: "Config (.aider.conf.yml)", Pattern: filepath.Join(home, ".aider.conf.yml")},
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
				{Label: "Workspace Dir", Pattern: "auto_gpt_workspace"},
				{Label: "SOUL.md", Pattern: "SOUL.md"},
				{Label: "USER.md", Pattern: "USER.md"},
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
				{Label: "SOUL.md", Pattern: "SOUL.md"},
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
			Name: "Goose", Ecosystem: "goose", Description: "Block's open-source AI dev agent", Icon: "🪿", CLIName: "goose",
			Expected: []models.ExpectedFile{
				{Label: "CLI Binary", Pattern: "goose"},
				{Label: "Config Dir", Pattern: filepath.Join(home, ".config", "goose")},
				{Label: "Config File", Pattern: filepath.Join(home, ".config", "goose", "config.yaml")},
				{Label: "SOUL.md", Pattern: filepath.Join(home, ".config", "goose", "SOUL.md")},
				{Label: "Profiles Dir", Pattern: filepath.Join(home, ".config", "goose", "profiles")},
				{Label: "Extensions Dir", Pattern: filepath.Join(home, ".config", "goose", "extensions")},
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
				{Label: "Hub Models (Windows)", Pattern: filepath.Join(home, ".lmstudio", "hub", "models")},
				{Label: "Hub Models (Linux)", Pattern: filepath.Join(home, ".cache", "lm-studio", "models")},
				{Label: "Hub Models (Linux alt)", Pattern: filepath.Join(home, ".local", "share", "lm-studio", "models")},
				{Label: "Hub Models (macOS)", Pattern: filepath.Join(home, ".lmstudio", "models")},
				{Label: "Config Dir", Pattern: lmStudioModels},
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
				{Label: "Tools Dir", Pattern: filepath.Join(home, ".open-webui", "tools")},
				{Label: "Functions Dir", Pattern: filepath.Join(home, ".open-webui", "functions")},
			},
		},
		{
			Name: "Continue", Ecosystem: "continue", Description: "VS Code AI coding assistant", Icon: "▶️", CLIName: "continue",
			Expected: []models.ExpectedFile{
				{Label: "Config File", Pattern: filepath.Join(home, ".continue", "config.json")},
				{Label: "Config TS", Pattern: filepath.Join(home, ".continue", "config.ts")},
				{Label: "Prompts Dir", Pattern: filepath.Join(home, ".continue", "prompts")},
				{Label: "Rules Dir", Pattern: filepath.Join(home, ".continue", "rules")},
				{Label: "Docs Dir", Pattern: filepath.Join(home, ".continue", "docs")},
			},
		},
	}
}

func lmStudioModelPath(home string) string {
	switch runtime.GOOS {
	case "windows":
		return filepath.Join(home, ".lmstudio", "hub", "models")
	case "darwin":
		return filepath.Join(home, ".lmstudio", "models")
	default:
		p := filepath.Join(home, ".cache", "lm-studio", "models")
		if _, err := os.Stat(p); err == nil {
			return p
		}
		return filepath.Join(home, ".local", "share", "lm-studio", "models")
	}
}

// ScanLMStudioModels scans the LM Studio hub directory structure:
// ~/.lmstudio/hub/models/{company}/{model-name}/*.gguf  (Windows)
func (ms *ModelScanner) ScanLMStudioModels() []models.AIModel {
	var result []models.AIModel
	home := ms.home
	roots := []string{
		filepath.Join(home, ".lmstudio", "hub", "models"),         // Windows new
		filepath.Join(home, ".lmstudio", "models"),                 // macOS
		filepath.Join(home, ".cache", "lm-studio", "models"),       // Linux
		filepath.Join(home, ".local", "share", "lm-studio", "models"), // Linux alt
	}
	for _, root := range roots {
		if _, err := os.Stat(root); os.IsNotExist(err) {
			continue
		}
		// Walk company/model/files structure
		companies, _ := os.ReadDir(root)
		for _, co := range companies {
			if !co.IsDir() {
				continue
			}
			modelDirs, _ := os.ReadDir(filepath.Join(root, co.Name()))
			for _, md := range modelDirs {
				if !md.IsDir() {
					continue
				}
				modelDir := filepath.Join(root, co.Name(), md.Name())
				files, _ := os.ReadDir(modelDir)
				for _, f := range files {
					if f.IsDir() {
						continue
					}
					ext := strings.ToLower(filepath.Ext(f.Name()))
					format, ok := modelExtensions[ext]
					if !ok {
						continue
					}
					info, err := f.Info()
					if err != nil || info.Size() < minModelSize {
						continue
					}
					result = append(result, models.AIModel{
						Name:       co.Name() + "/" + md.Name(),
						Tag:        f.Name(),
						SizeBytes:  info.Size(),
						Provider:   "LM Studio",
						FilePath:   filepath.Join(modelDir, f.Name()),
						Extension:  ext,
						Format:     format,
						IsOrphaned: false,
					})
				}
			}
		}
	}
	return result
}

// ScanMemoryFiles does a targeted search for soul/memory markdown files
// across the entire home directory shallowly, plus known agent dirs deeply.
func ScanMemoryFiles(home string) []models.AgentFile {
	var files []models.AgentFile
	seen := map[string]bool{}

	// Deep scan of known agent config dirs
	agentDirs := []string{
		filepath.Join(home, ".claude"),
		filepath.Join(home, ".openclaw"),
		filepath.Join(home, ".hermes"),
		filepath.Join(home, ".openinterpreter"),
		filepath.Join(home, ".continue"),
		filepath.Join(home, ".config", "goose"),
		filepath.Join(home, ".aider.chat.history.md"),
	}

	for _, dir := range agentDirs {
		filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
			if err != nil || d.IsDir() || seen[path] {
				return nil
			}
			name := d.Name()
			ext := strings.ToLower(filepath.Ext(name))
			baseName := strings.ToUpper(strings.TrimSuffix(name, filepath.Ext(name)))

			isMemory := ext == ".soul" || ext == ".persona" || ext == ".memory"
			for _, mf := range memoryFileNames {
				if strings.EqualFold(name, mf) {
					isMemory = true
					break
				}
			}
			_ = baseName
			if !isMemory {
				return nil
			}
			seen[path] = true
			info, _ := d.Info()
			var sz int64
			if info != nil {
				sz = info.Size()
			}
			files = append(files, models.AgentFile{
				Type: classifyMemoryFile(name), Name: name, Path: path, SizeBytes: sz,
			})
			return nil
		})
	}

	// Shallow home scan for loose memory files
	entries, _ := os.ReadDir(home)
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		for _, mf := range memoryFileNames {
			if strings.EqualFold(e.Name(), mf) {
				path := filepath.Join(home, e.Name())
				if seen[path] {
					continue
				}
				seen[path] = true
				info, _ := e.Info()
				var sz int64
				if info != nil {
					sz = info.Size()
				}
				files = append(files, models.AgentFile{
					Type: classifyMemoryFile(e.Name()), Name: e.Name(), Path: path, SizeBytes: sz,
				})
			}
		}
	}

	return files
}

func classifyMemoryFile(name string) string {
	upper := strings.ToUpper(name)
	switch {
	case strings.Contains(upper, "SOUL"):
		return "soul"
	case strings.Contains(upper, "USER") || strings.Contains(upper, "PROFILE") || strings.Contains(upper, "PERSONA") || strings.Contains(upper, "IDENTITY"):
		return "persona"
	case strings.Contains(upper, "MEMORY") || strings.Contains(upper, "SUMMARY") || strings.Contains(upper, "HISTORY") || strings.Contains(upper, "NOTES") || strings.Contains(upper, "CONTEXT"):
		return "memory"
	case strings.Contains(upper, "AGENT") || strings.Contains(upper, "INSTRUCTIONS") || strings.Contains(upper, "RULES") || strings.Contains(upper, "SYSTEM"):
		return "agent-def"
	case strings.Contains(upper, "SKILL"):
		return "skill"
	case strings.Contains(upper, "KNOWLEDGE"):
		return "knowledge"
	default:
		return "memory"
	}
}

func resolveExpectedFile(home, cwd, pattern, cliName string) (found bool, path string) {
	if filepath.IsAbs(pattern) {
		if _, err := os.Stat(pattern); err == nil {
			return true, pattern
		}
		return false, ""
	}
	// CLI binary
	if !strings.Contains(pattern, string(os.PathSeparator)) && !strings.Contains(pattern, ".") {
		if p, err := exec.LookPath(pattern); err == nil {
			return true, p
		}
	}
	// Relative to cwd
	if p := filepath.Join(cwd, pattern); fileExists(p) {
		return true, p
	}
	// Relative to home
	if p := filepath.Join(home, pattern); fileExists(p) {
		return true, p
	}
	return false, ""
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func BuildAgentCards(home string, progressCh chan<- string) []models.AgentCard {
	cwd, _ := os.Getwd()
	defs := getAgentDefs(home)
	var cards []models.AgentCard

	// Pre-scan memory files once and attach to relevant cards
	memFiles := ScanMemoryFiles(home)

	for _, def := range defs {
		if progressCh != nil {
			progressCh <- "Profiling " + def.Name + " ecosystem..."
		}

		card := models.AgentCard{
			Name: def.Name, Ecosystem: def.Ecosystem,
			Description: def.Description, Icon: def.Icon, Status: "not_found",
		}

		foundCount := 0
		for _, ef := range def.Expected {
			found, resolvedPath := resolveExpectedFile(home, cwd, ef.Pattern, def.CLIName)
			ef.Found = found
			ef.Path = resolvedPath
			if found {
				foundCount++
			}
			card.ExpectedFiles = append(card.ExpectedFiles, ef)
		}

		if foundCount == len(def.Expected) && len(def.Expected) > 0 {
			card.Status = "detected"
		} else if foundCount > 0 {
			card.Status = "partial"
		}

		// Attach ecosystem-specific files + memory files
		card.Files = append(scanAgentFiles(home, def.Ecosystem), filterMemoryFilesForEcosystem(memFiles, def.Ecosystem, home)...)
		cards = append(cards, card)
	}

	return cards
}

func filterMemoryFilesForEcosystem(all []models.AgentFile, ecosystem, home string) []models.AgentFile {
	var out []models.AgentFile
	ecoDir := filepath.Join(home, "."+ecosystem)
	for _, f := range all {
		if strings.HasPrefix(f.Path, ecoDir) {
			out = append(out, f)
		}
	}
	return out
}

func scanAgentFiles(home, ecosystem string) []models.AgentFile {
	var files []models.AgentFile
	ecosystemDirs := map[string][]string{
		"claude":          {filepath.Join(home, ".claude")},
		"openclaw":        {filepath.Join(home, ".openclaw")},
		"hermes":          {filepath.Join(home, ".hermes")},
		"openinterpreter": {filepath.Join(home, ".openinterpreter")},
		"aider":           {home},
		"crewai":          {".", "crewai"},
		"langchain":       {filepath.Join(home, ".langchain"), "."},
		"goose":           {filepath.Join(home, ".config", "goose")},
		"ollama":          {filepath.Join(home, ".ollama")},
		"continue":        {filepath.Join(home, ".continue")},
		"lmstudio":        {filepath.Join(home, ".lmstudio"), filepath.Join(home, ".cache", "lm-studio")},
		"openwebui":       {filepath.Join(home, ".open-webui")},
	}

	fileTypeMap := map[string]string{
		".yaml": "config", ".yml": "config", ".json": "config",
		".md": "memory", ".toml": "config", ".env": "config",
		".db": "memory", ".gguf": "model",
		".soul": "soul", ".skill": "skill", ".tool": "tool", ".persona": "persona",
	}

	dirs, ok := ecosystemDirs[ecosystem]
	if !ok {
		return files
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
			// Classify special memory files
			for _, mf := range memoryFileNames {
				if strings.EqualFold(entry.Name(), mf) {
					ftype = classifyMemoryFile(entry.Name())
					break
				}
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
