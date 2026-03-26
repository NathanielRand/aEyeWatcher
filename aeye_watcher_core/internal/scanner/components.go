package scanner

import (
	"fmt"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/aeye/watcher-core/internal/models"
)

type ComponentDef struct {
	Name        string
	Type        string
	Description string
	CLIName     string
	Port        string
	Endpoint    string
	AltPorts    []string
}

var componentDefs = []ComponentDef{
	// Core Providers
	{Name: "Ollama", Type: "Provider", Description: "Local LLM server", CLIName: "ollama", Port: "11434", Endpoint: "/"},
	{Name: "LM Studio", Type: "Provider", Description: "Desktop LLM runner (lms CLI)", CLIName: "lms", Port: "1234", Endpoint: "/v1/models"},
	{Name: "LocalAI", Type: "Provider", Description: "OpenAI-compatible local server", CLIName: "local-ai", Port: "8080", Endpoint: "/v1/models"},
	{Name: "llama-server", Type: "Provider", Description: "llama.cpp HTTP server", CLIName: "llama-server", Port: "8000", Endpoint: "/health"},
	{Name: "koboldcpp", Type: "Provider", Description: "KoboldCPP inference server", CLIName: "koboldcpp", Port: "5001", Endpoint: "/api/v1/model"},
	{Name: "vLLM", Type: "Provider", Description: "High-throughput LLM serving", CLIName: "vllm", Port: "8000", Endpoint: "/v1/models"},
	{Name: "text-generation-webui", Type: "Provider", Description: "Oobabooga WebUI server", CLIName: "text-generation-webui", Port: "7860", Endpoint: "/"},
	{Name: "Jan", Type: "Provider", Description: "Jan.ai local server", CLIName: "jan", Port: "1337", Endpoint: "/"},
	{Name: "GPT4All", Type: "Provider", Description: "GPT4All local server", CLIName: "gpt4all", Port: "4891", Endpoint: "/v1/models"},
	{Name: "Llamafile", Type: "Provider", Description: "Self-contained model server", CLIName: "llamafile", Port: "8080", Endpoint: "/v1/models"},

	// Agents
	{Name: "Claude Code", Type: "Agent", Description: "Anthropic agentic coding assistant", CLIName: "claude"},
	{Name: "OpenInterpreter", Type: "Agent", Description: "Open-source code interpreter agent", CLIName: "interpreter"},
	{Name: "Aider", Type: "Agent", Description: "AI pair programmer in your terminal", CLIName: "aider"},
	{Name: "AutoGPT", Type: "Agent", Description: "Autonomous GPT agent framework", CLIName: "autogpt"},
	{Name: "CrewAI", Type: "Agent", Description: "Multi-agent orchestration", CLIName: "crewai"},
	{Name: "LangChain", Type: "Agent", Description: "LLM application framework CLI", CLIName: "langchain"},
	{Name: "SWE-agent", Type: "Agent", Description: "Software engineering agent", CLIName: "sweagent"},
	{Name: "Devin", Type: "Agent", Description: "Devin CLI agent", CLIName: "devin"},
	{Name: "Goose", Type: "Agent", Description: "Block open-source AI agent", CLIName: "goose"},
	{Name: "Amp", Type: "Agent", Description: "Sourcegraph coding agent", CLIName: "amp"},
	// OpenClaw - also opens a gateway on boot
	{Name: "OpenClaw", Type: "Agent", Description: "OpenClaw multi-agent framework (opens gateway on boot)", CLIName: "openclaw", Port: "7331", Endpoint: "/health", AltPorts: []string{"7332", "7333"}},
	// Hermes - agent + exposes HTTP gateway
	{Name: "Hermes", Type: "Agent", Description: "Hermes agent runtime (exposes HTTP gateway)", CLIName: "hermes", Port: "8765", Endpoint: "/health", AltPorts: []string{"8766", "8767"}},

	// Agent Skills
	{Name: "Continue", Type: "Agent Skill", Description: "VS Code AI extension", CLIName: "continue"},
	{Name: "Copilot CLI", Type: "Agent Skill", Description: "GitHub Copilot CLI", CLIName: "gh-copilot"},
	{Name: "Cursor", Type: "Agent Skill", Description: "Cursor AI IDE", CLIName: "cursor"},
	{Name: "Windsurf", Type: "Agent Skill", Description: "Codeium Windsurf IDE", CLIName: "windsurf"},

	// Gateways & Routers
	{Name: "LiteLLM", Type: "Gateway", Description: "OpenAI-compatible proxy/router", CLIName: "litellm", Port: "4000", Endpoint: "/health"},
	{Name: "OpenRouter", Type: "Gateway", Description: "Cloud model routing service", CLIName: "openrouter"},
	{Name: "Portkey", Type: "Gateway", Description: "AI gateway with observability", CLIName: "portkey", Port: "8787", Endpoint: "/v1/health"},
	{Name: "OpenClaw Gateway", Type: "Gateway", Description: "OpenClaw boot gateway (auto-started by agent)", Port: "7331", Endpoint: "/gateway/health"},
	{Name: "Hermes Gateway", Type: "Gateway", Description: "Hermes agent gateway", Port: "8765", Endpoint: "/gateway"},
	{Name: "HelixML", Type: "Gateway", Description: "HelixML AI gateway", CLIName: "helix", Port: "8090", Endpoint: "/healthz"},

	// Vector DBs
	{Name: "ChromaDB", Type: "Vector DB", Description: "Embedding database", CLIName: "chroma", Port: "8000", Endpoint: "/api/v1"},
	{Name: "Qdrant", Type: "Vector DB", Description: "Vector search engine", CLIName: "qdrant", Port: "6333", Endpoint: "/"},
	{Name: "Weaviate", Type: "Vector DB", Description: "Vector database", CLIName: "weaviate", Port: "8080", Endpoint: "/v1/.well-known/ready"},
	{Name: "Milvus", Type: "Vector DB", Description: "Open-source vector DB", CLIName: "milvus", Port: "19530", Endpoint: "/"},
	{Name: "LanceDB", Type: "Vector DB", Description: "Embedded vector DB", CLIName: "lancedb"},
	{Name: "pgvector", Type: "Vector DB", Description: "Postgres vector extension", CLIName: "psql"},

	// GUI Clients
	{Name: "Open WebUI", Type: "GUI Client", Description: "Web UI for local models", CLIName: "open-webui", Port: "3000", Endpoint: "/"},
	{Name: "AnythingLLM", Type: "GUI Client", Description: "All-in-one LLM desktop client", CLIName: "anythingllm", Port: "3001", Endpoint: "/"},
	{Name: "Msty", Type: "GUI Client", Description: "Desktop AI chat client", CLIName: "msty"},
	{Name: "LibreChat", Type: "GUI Client", Description: "Open-source chat UI", CLIName: "librechat", Port: "3080", Endpoint: "/"},

	// TTS / Audio
	{Name: "Whisper.cpp", Type: "TTS / Audio", Description: "Fast speech recognition", CLIName: "whisper-cpp"},
	{Name: "Whisper", Type: "TTS / Audio", Description: "OpenAI Whisper", CLIName: "whisper"},
	{Name: "Piper", Type: "TTS / Audio", Description: "Neural TTS", CLIName: "piper"},
	{Name: "Bark", Type: "TTS / Audio", Description: "Text-to-audio model", CLIName: "bark"},
	{Name: "Coqui TTS", Type: "TTS / Audio", Description: "Coqui TTS server", CLIName: "tts", Port: "5002", Endpoint: "/"},

	// Observability
	{Name: "Langfuse", Type: "Observability", Description: "LLM observability platform", CLIName: "langfuse", Port: "3000", Endpoint: "/api/health"},
	{Name: "Phoenix", Type: "Observability", Description: "Arize Phoenix tracing", CLIName: "phoenix", Port: "6006", Endpoint: "/"},
	{Name: "Helicone", Type: "Observability", Description: "LLM observability proxy", CLIName: "helicone"},
	{Name: "Weights & Biases", Type: "Observability", Description: "ML experiment tracking", CLIName: "wandb"},
}

func ScanComponents(progressCh chan<- string) []models.Component {
	var result []models.Component
	client := http.Client{Timeout: 800 * time.Millisecond}

	for _, def := range componentDefs {
		if progressCh != nil {
			progressCh <- fmt.Sprintf("Checking %s...", def.Name)
		}

		comp := models.Component{
			Name: def.Name, Type: def.Type,
			Description: def.Description, Status: models.StatusNotFound,
		}

		if def.CLIName != "" {
			if path, err := exec.LookPath(def.CLIName); err == nil {
				comp.Status = models.StatusInstalled
				comp.Path = path
				comp.Version = getVersion(def.CLIName)
			}
		}

		if def.Port != "" && def.Endpoint != "" {
			if probePort(client, def.Port, def.Endpoint) {
				comp.Status = models.StatusRunning
				comp.Port = def.Port
				comp.URL = fmt.Sprintf("http://localhost:%s%s", def.Port, def.Endpoint)
			}
		}

		if comp.Status != models.StatusRunning {
			for _, alt := range def.AltPorts {
				if probePort(client, alt, def.Endpoint) {
					comp.Status = models.StatusRunning
					comp.Port = alt
					comp.URL = fmt.Sprintf("http://localhost:%s%s", alt, def.Endpoint)
					break
				}
			}
		}

		result = append(result, comp)
	}

	return result
}

func probePort(client http.Client, port, endpoint string) bool {
	url := fmt.Sprintf("http://localhost:%s%s", port, endpoint)
	resp, err := client.Get(url)
	if err != nil {
		return false
	}
	resp.Body.Close()
	return resp.StatusCode < 500
}

func getVersion(cli string) string {
	cmds := [][]string{{cli, "--version"}, {cli, "version"}, {cli, "-v"}}
	for _, args := range cmds {
		out, err := exec.Command(args[0], args[1:]...).Output()
		if err == nil {
			v := strings.TrimSpace(string(out))
			if len(v) > 0 && len(v) < 100 {
				return strings.Split(v, "\n")[0]
			}
		}
	}
	return ""
}
