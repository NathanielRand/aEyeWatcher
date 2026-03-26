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
	{Name: "jan", Type: "Provider", Description: "Jan.ai local server", CLIName: "jan", Port: "1337", Endpoint: "/"},

	// Agents
	{Name: "Claude Code", Type: "Agent", Description: "Anthropic CLI coding agent", CLIName: "claude"},
	{Name: "OpenInterpreter", Type: "Agent", Description: "Open interpreter agent", CLIName: "interpreter"},
	{Name: "Aider", Type: "Agent", Description: "AI pair programmer", CLIName: "aider"},
	{Name: "AutoGPT", Type: "Agent", Description: "Autonomous GPT agent", CLIName: "autogpt"},
	{Name: "CrewAI", Type: "Agent", Description: "Multi-agent orchestration", CLIName: "crewai"},
	{Name: "LangChain", Type: "Agent", Description: "LangChain CLI", CLIName: "langchain"},
	{Name: "Hermes", Type: "Agent", Description: "Hermes agent framework", CLIName: "hermes"},
	{Name: "OpenClaw", Type: "Agent", Description: "OpenClaw agent", CLIName: "openclaw"},
	{Name: "SWE-agent", Type: "Agent", Description: "Software engineering agent", CLIName: "sweagent"},
	{Name: "Devin", Type: "Agent", Description: "Devin CLI agent", CLIName: "devin"},
	{Name: "Continue", Type: "Agent Skill", Description: "VS Code AI extension", CLIName: "continue"},
	{Name: "Copilot CLI", Type: "Agent Skill", Description: "GitHub Copilot CLI", CLIName: "gh-copilot"},

	// Gateways & Routers
	{Name: "LiteLLM", Type: "Gateway", Description: "OpenAI-compatible proxy/router", CLIName: "litellm", Port: "4000", Endpoint: "/health"},
	{Name: "OpenRouter", Type: "Gateway", Description: "Model routing service", CLIName: "openrouter"},
	{Name: "Portkey", Type: "Gateway", Description: "AI gateway", CLIName: "portkey"},

	// Vector DBs
	{Name: "ChromaDB", Type: "Vector DB", Description: "Embedding database", CLIName: "chroma", Port: "8000", Endpoint: "/api/v1"},
	{Name: "Qdrant", Type: "Vector DB", Description: "Vector search engine", CLIName: "qdrant", Port: "6333", Endpoint: "/"},
	{Name: "Weaviate", Type: "Vector DB", Description: "Vector database", CLIName: "weaviate", Port: "8080", Endpoint: "/v1/.well-known/ready"},
	{Name: "Milvus", Type: "Vector DB", Description: "Open-source vector DB", CLIName: "milvus", Port: "19530", Endpoint: "/"},
	{Name: "Pinecone", Type: "Vector DB", Description: "Managed vector DB", CLIName: "pinecone"},
	{Name: "pgvector", Type: "Vector DB", Description: "Postgres vector extension", CLIName: "psql"},

	// GUI Clients
	{Name: "Open WebUI", Type: "GUI Client", Description: "Web UI for local models", CLIName: "open-webui", Port: "3000", Endpoint: "/"},
	{Name: "AnythingLLM", Type: "GUI Client", Description: "All-in-one LLM desktop client", CLIName: "anythingllm"},
	{Name: "Msty", Type: "GUI Client", Description: "Desktop AI chat client", CLIName: "msty"},

	// TTS / Audio
	{Name: "Whisper.cpp", Type: "TTS / Audio", Description: "Speech recognition", CLIName: "whisper-cpp"},
	{Name: "whisper", Type: "TTS / Audio", Description: "OpenAI Whisper", CLIName: "whisper"},
	{Name: "Piper", Type: "TTS / Audio", Description: "Neural TTS", CLIName: "piper"},
	{Name: "Bark", Type: "TTS / Audio", Description: "Text-to-audio model", CLIName: "bark"},

	// Monitoring
	{Name: "Langfuse", Type: "Observability", Description: "LLM observability", CLIName: "langfuse", Port: "3000", Endpoint: "/api/health"},
	{Name: "Phoenix", Type: "Observability", Description: "Arize Phoenix tracing", CLIName: "phoenix", Port: "6006", Endpoint: "/"},
}

func ScanComponents(progressCh chan<- string) []models.Component {
	var result []models.Component
	client := http.Client{Timeout: 800 * time.Millisecond}

	for _, def := range componentDefs {
		if progressCh != nil {
			progressCh <- fmt.Sprintf("Checking %s...", def.Name)
		}

		comp := models.Component{
			Name:        def.Name,
			Type:        def.Type,
			Description: def.Description,
			Status:      models.StatusNotFound,
		}

		// Check CLI
		if def.CLIName != "" {
			if path, err := exec.LookPath(def.CLIName); err == nil {
				comp.Status = models.StatusInstalled
				comp.Path = path
				// Try to get version
				comp.Version = getVersion(def.CLIName)
			}
		}

		// Check service on port
		if def.Port != "" && def.Endpoint != "" {
			url := fmt.Sprintf("http://localhost:%s%s", def.Port, def.Endpoint)
			resp, err := client.Get(url)
			if err == nil && resp.StatusCode < 500 {
				comp.Status = models.StatusRunning
				comp.Port = def.Port
				comp.URL = url
				if resp != nil {
					resp.Body.Close()
				}
			}
		}

		result = append(result, comp)
	}

	return result
}

func getVersion(cli string) string {
	cmds := [][]string{
		{cli, "--version"},
		{cli, "version"},
		{cli, "-v"},
	}
	for _, args := range cmds {
		out, err := exec.Command(args[0], args[1:]...).Output()
		if err == nil {
			v := strings.TrimSpace(string(out))
			if len(v) > 0 && len(v) < 80 {
				// Extract just the version line
				lines := strings.Split(v, "\n")
				return lines[0]
			}
		}
	}
	return ""
}
