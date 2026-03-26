# ◉ aEye Watcher v2

> *"Every model leaves a trace. Every agent leaves a trail. We watch them all."*

**aEye Watcher** is a full-stack AI ecosystem monitor — a living command center that scans your entire system for AI providers, models, agents, MCP servers, ACP agents, vector databases, and more. Built with a vintage detective / grainy comic-book command center aesthetic.

---

## Architecture

```
aeye_watcher_core/     ← Go backend API + WebSocket scan engine
aeye_watcher_ui/       ← SvelteKit 5 + Tailwind frontend
```

The two are fully decoupled. The core exposes REST + WebSocket endpoints; any frontend can consume it.

---

## aeye_watcher_core (Go)

### Features
- **OS Detection** — Linux (distro/version), macOS, Windows; hostname, arch, container detection
- **GPU Monitoring** — NVIDIA (nvidia-smi), AMD (ROCm), Apple Silicon; live VRAM usage
- **Deep Model Scanner** — hunts for `.gguf`, `.safetensors`, `.onnx`, `.pt`, `.pth`, `.bin`, `.mlmodel`, `.tflite`, `.engine`, `.llamafile` across your entire filesystem. Tags files outside known provider dirs as **orphaned**
- **Provider Parsing** — Ollama manifests, HuggingFace hub cache, LM Studio cache
- **Component Detector** — 30+ components: Ollama, LM Studio, LocalAI, vLLM, llama-server, KoboldCPP, Jan, Aider, Claude Code, AutoGPT, CrewAI, LangChain, LiteLLM, ChromaDB, Qdrant, Weaviate, Open WebUI, Whisper, Piper, Langfuse, and more
- **MCP Server Scanner** — parses Claude Desktop, Cursor, Windsurf, Continue, Zed configs; discovers `mcpServers` blocks; probes ports; walks `~` for any `mcp*.json`
- **ACP Agent Scanner** — probes `/.well-known/agent.json` on standard agent ports
- **Agent Ecosystem Cards** — profiles 10 agent ecosystems (Claude Code, Aider, CrewAI, Ollama, LM Studio, OpenInterpreter, AutoGPT, LangChain, Open WebUI, Continue) with expected file checklists
- **Disk Scanner** — maps all real volumes with usage stats; cross-platform (Linux `df`, Windows `wmic`)
- **WebSocket Scan Progress** — phase-by-phase live events with percent completion, checkpoints, and found items
- **Safe Mode Deletion** — logs cleanup requests without touching files; opt-in via `AEYE_UNSAFE_DELETE=1`

### API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| `WS` | `/ws/scan` | Live scanning WebSocket — streams `ScanEvent` JSON |
| `GET` | `/api/gpu` | Current GPU stats (pollable) |
| `GET` | `/api/os` | OS/system info |
| `GET` | `/api/scan/cached` | Last scan result (JSON) |
| `POST` | `/api/delete` | Request model cleanup (safe mode) |
| `GET` | `/health` | Health check |

### WebSocket Event Schema

```json
{
  "type": "checkpoint | progress | found | complete | error",
  "phase": "system | gpu | disk | components | mcp | acp | models | agents | complete",
  "message": "Human-readable status string",
  "percent": 0-100,
  "detail": { ... }  // present on 'found' and 'complete' events
}
```

The `complete` event's `detail` field is the full `ScanResult` object.

### Setup & Run

```bash
cd aeye_watcher_core

# Install dependencies
go mod tidy

# Run (default port 8080)
go run main.go

# Custom port
go run main.go -port 9090

# Build binary
go build -o aeye-core main.go
./aeye-core
```

**Requirements:** Go 1.22+. Optional: `nvidia-smi` for GPU monitoring.

---

## aeye_watcher_ui (SvelteKit)

### Features
- **Svelte 5 runes mode** — `$state`, `$derived`, `$effect` throughout
- **5 themes** — Obsidian (purple), Clay (orange), Forest (green), Slate (blue), Amber (yellow) — all muted/dark with grainy overlay
- **Vintage detective aesthetic** — Special Elite + Courier Prime + Rajdhani fonts, grain SVG overlay, animated glowing dots, comic-book-style iconography
- **Live WebSocket scan** with animated progress bar, phase indicators, scrolling live feed
- **Command Center Dashboard** — stat cards, live provider chips, disk overview, activity feed
- **Providers View** — grouped by type (Providers, Gateways, Vector DBs, GUIs, TTS, Observability)
- **Agent Ecosystem Cards** — click-through to full character card with expected-file checklist
- **MCP / ACP View** — all discovered MCP configs and live ACP agents
- **Model Cache** — searchable/filterable model list with orphan tagging, format badges, size, cleanup
- **Storage View** — full disk cards with usage bars
- **Activity Log** — timestamped live event feed

### Setup & Run

```bash
cd aeye_watcher_ui

# Install deps
npm install

# Dev server (proxies /api and /ws to localhost:8080)
npm run dev

# Build for production
npm run build
npm run preview
```

**Requirements:** Node 18+. The Go core must be running on port 8080 (or update `CORE_URL` in `src/lib/stores/app.svelte.ts`).

---

## Scan Phases

| Phase | What Happens |
|-------|-------------|
| `system` | OS detection, hostname, arch |
| `gpu` | nvidia-smi / rocm-smi / sysctl probe |
| `disk` | `df` / `wmic` volume enumeration |
| `components` | 30+ CLI + port checks for providers, agents, gateways, etc. |
| `mcp` | Config file walking + port probing for MCP servers |
| `acp` | `/.well-known/agent.json` probing on agent ports |
| `models` | Ollama manifests → HF cache → deep filesystem walk |
| `agents` | Agent ecosystem profiling & file checklists |
| `complete` | Full `ScanResult` emitted on WebSocket |

---

## Orphaned Model Detection

A model file is tagged **orphaned** if it exists outside all known provider directories:
- `~/.ollama/models`
- `~/.cache/huggingface/hub`
- `~/.cache/lm-studio/models`
- `~/.local/share/lm-studio/models`
- `~/Library/Application Support/LM Studio` (macOS)
- `koboldcpp`, `llama.cpp`, `text-generation-webui/models`, `oobabooga`

Orphaned files are flagged for cleanup in the Models view.

---

## Theme System

Themes are applied via `data-theme` on `<html>` and CSS variables:

| Theme | Accent | Vibe |
|-------|--------|------|
| `obsidian` | `#c084fc` | Purple — default dark ops |
| `clay` | `#fb923c` | Warm orange — dusty detective |
| `forest` | `#4ade80` | Green — field operative |
| `slate` | `#7dd3fc` | Steel blue — cold intelligence |
| `amber` | `#fbbf24` | Golden — noir lamp light |

Selection persists to `localStorage`.

---

## Enabling Real Deletion

By default the delete API is in **SAFE MODE** — it logs requests but does not touch files. To enable:

```bash
AEYE_UNSAFE_DELETE=1 go run main.go
```

Then uncomment the `os.Remove` block in `internal/server/server.go`.

---

## Extending the Scanner

**Add a new component** → append to `componentDefs` in `internal/scanner/components.go`

**Add a new model extension** → append to `modelExtensions` in `internal/scanner/models.go`

**Add a new agent ecosystem** → append to `getAgentDefs()` in `internal/scanner/agents.go`

**Add a new MCP config path** → append to `getMCPConfigPaths()` in `internal/scanner/mcp_acp.go`

---

## Project Structure

```
aeye_watcher_core/
├── main.go
├── go.mod
├── go.sum
└── internal/
    ├── models/
    │   └── types.go          # All shared data types
    ├── os_detect/
    │   └── detect.go         # OS/arch/GPU detection
    ├── disk/
    │   └── disk.go           # Cross-platform disk scanner
    ├── scanner/
    │   ├── components.go     # Provider/agent/gateway detection
    │   ├── models.go         # Model file hunting (Ollama, HF, deep scan)
    │   ├── mcp_acp.go        # MCP server + ACP agent discovery
    │   ├── agents.go         # Agent ecosystem card builder
    │   └── gpu.go            # GPU stats (NVIDIA/AMD/Apple)
    └── server/
        └── server.go         # HTTP + WebSocket server

aeye_watcher_ui/
├── package.json
├── svelte.config.js
├── vite.config.ts
├── tailwind.config.js
├── tsconfig.json
└── src/
    ├── app.html
    ├── app.css
    ├── lib/
    │   ├── types.ts              # TypeScript mirrors of Go types
    │   └── stores/
    │       └── app.svelte.ts     # All reactive state + API calls
    └── routes/
        ├── +layout.svelte
        └── +page.svelte          # Full command center UI
```
