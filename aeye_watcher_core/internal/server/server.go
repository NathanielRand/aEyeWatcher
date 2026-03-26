package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"time"

	"github.com/gorilla/websocket"

	"github.com/aeye/watcher-core/internal/disk"
	"github.com/aeye/watcher-core/internal/models"
	"github.com/aeye/watcher-core/internal/os_detect"
	"github.com/aeye/watcher-core/internal/scanner"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Server struct {
	osInfo     os_detect.OSInfo
	cachedScan *models.ScanResult
	scanTime   time.Time
}

func New() *Server {
	return &Server{
		osInfo: os_detect.Detect(),
	}
}

func (s *Server) withCORS(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		h(w, r)
	}
}

// ScanWebSocket handles the live scanning experience via WebSocket
func (s *Server) ScanWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WS upgrade error:", err)
		return
	}
	defer conn.Close()

	sendEvent := func(evt models.ScanEvent) {
		data, _ := json.Marshal(evt)
		conn.WriteMessage(websocket.TextMessage, data)
	}

	home, _ := os.UserHomeDir()

	// Phase 1: OS & System
	sendEvent(models.ScanEvent{Type: models.EventCheckpoint, Phase: "system", Message: "Initializing system profile...", Percent: 2})
	time.Sleep(200 * time.Millisecond)
	sendEvent(models.ScanEvent{Type: models.EventFound, Phase: "system", Message: fmt.Sprintf("OS detected: %s %s (%s)", s.osInfo.Distro, s.osInfo.Version, s.osInfo.Arch), Percent: 5})

	// Phase 2: GPU
	sendEvent(models.ScanEvent{Type: models.EventCheckpoint, Phase: "gpu", Message: "Interrogating GPU...", Percent: 8})
	gpu := scanner.ScanGPU()
	if gpu.Error == "" {
		sendEvent(models.ScanEvent{Type: models.EventFound, Phase: "gpu", Message: fmt.Sprintf("GPU: %s (%d MB VRAM)", gpu.Name, gpu.TotalMB), Percent: 12, Detail: gpu})
	} else {
		sendEvent(models.ScanEvent{Type: models.EventFound, Phase: "gpu", Message: "No GPU monitoring available", Percent: 12})
	}

	// Phase 3: Disk
	sendEvent(models.ScanEvent{Type: models.EventCheckpoint, Phase: "disk", Message: "Mapping storage volumes...", Percent: 15})
	disks := disk.ScanDisks()
	sendEvent(models.ScanEvent{Type: models.EventFound, Phase: "disk", Message: fmt.Sprintf("Found %d disk volume(s)", len(disks)), Percent: 20, Detail: disks})

	// Phase 4: Components
	sendEvent(models.ScanEvent{Type: models.EventCheckpoint, Phase: "components", Message: "Scanning AI providers, agents, gateways...", Percent: 22})
	progressCh := make(chan string, 50)
	var components []models.Component
	go func() {
		components = scanner.ScanComponents(progressCh)
		close(progressCh)
	}()

	basePercent := 22
	for msg := range progressCh {
		sendEvent(models.ScanEvent{Type: models.EventProgress, Phase: "components", Message: msg, Percent: basePercent})
		if basePercent < 45 {
			basePercent++
		}
	}
	running := 0
	for _, c := range components {
		if c.Status == models.StatusRunning {
			running++
		}
	}
	sendEvent(models.ScanEvent{Type: models.EventFound, Phase: "components", Message: fmt.Sprintf("Detected %d components (%d running)", len(components), running), Percent: 45})

	// Phase 5: MCP Servers
	sendEvent(models.ScanEvent{Type: models.EventCheckpoint, Phase: "mcp", Message: "Hunting MCP server configurations...", Percent: 48})
	mcpCh := make(chan string, 20)
	var mcpServers []models.MCPServer
	go func() {
		mcpServers = scanner.ScanMCPServers(home, mcpCh)
		close(mcpCh)
	}()
	for msg := range mcpCh {
		sendEvent(models.ScanEvent{Type: models.EventProgress, Phase: "mcp", Message: msg, Percent: 50})
	}
	sendEvent(models.ScanEvent{Type: models.EventFound, Phase: "mcp", Message: fmt.Sprintf("Found %d MCP server configuration(s)", len(mcpServers)), Percent: 55})

	// Phase 6: ACP Agents
	sendEvent(models.ScanEvent{Type: models.EventCheckpoint, Phase: "acp", Message: "Probing ACP agent endpoints...", Percent: 57})
	acpAgents := scanner.ScanACPAgents(nil)
	sendEvent(models.ScanEvent{Type: models.EventFound, Phase: "acp", Message: fmt.Sprintf("Found %d ACP agent(s)", len(acpAgents)), Percent: 60})

	// Phase 7: Models - Ollama manifests
	sendEvent(models.ScanEvent{Type: models.EventCheckpoint, Phase: "models", Message: "Parsing Ollama model manifests...", Percent: 62})
	ms := scanner.NewModelScanner(home, nil)
	ollamaModels := ms.ScanOllamaManifests()
	sendEvent(models.ScanEvent{Type: models.EventFound, Phase: "models", Message: fmt.Sprintf("Ollama: %d model(s)", len(ollamaModels)), Percent: 65})

	// HuggingFace cache
	sendEvent(models.ScanEvent{Type: models.EventCheckpoint, Phase: "models", Message: "Parsing HuggingFace cache...", Percent: 67})
	hfModels := ms.ScanHuggingFaceCache()
	sendEvent(models.ScanEvent{Type: models.EventFound, Phase: "models", Message: fmt.Sprintf("HuggingFace: %d model(s)", len(hfModels)), Percent: 70})

	// Deep file scan
	sendEvent(models.ScanEvent{Type: models.EventCheckpoint, Phase: "models", Message: "Deep scanning filesystem for AI model files...", Percent: 72})
	deepCh := make(chan string, 100)
	var deepModels []models.AIModel
	go func() {
		deepModels = ms.DeepScan(deepCh)
		close(deepCh)
	}()
	deepPercent := 72
	msgCount := 0
	for msg := range deepCh {
		msgCount++
		if msgCount%5 == 0 { // throttle updates
			sendEvent(models.ScanEvent{Type: models.EventProgress, Phase: "models", Message: msg, Percent: deepPercent})
			if deepPercent < 85 {
				deepPercent++
			}
		}
	}

	allModels := scanner.DeduplicateModels(append(append(ollamaModels, hfModels...), deepModels...))
	orphaned := 0
	for _, m := range allModels {
		if m.IsOrphaned {
			orphaned++
		}
	}
	sendEvent(models.ScanEvent{Type: models.EventFound, Phase: "models", Message: fmt.Sprintf("Total models: %d (%d orphaned)", len(allModels), orphaned), Percent: 86})

	// Phase 8: Agent Cards
	sendEvent(models.ScanEvent{Type: models.EventCheckpoint, Phase: "agents", Message: "Building agent ecosystem profiles...", Percent: 88})
	agentCh := make(chan string, 30)
	var agentCards []models.AgentCard
	go func() {
		agentCards = scanner.BuildAgentCards(home, agentCh)
		close(agentCh)
	}()
	for msg := range agentCh {
		sendEvent(models.ScanEvent{Type: models.EventProgress, Phase: "agents", Message: msg, Percent: 90})
	}
	sendEvent(models.ScanEvent{Type: models.EventFound, Phase: "agents", Message: fmt.Sprintf("Profiled %d agent ecosystems", len(agentCards)), Percent: 94})

	// Sort models by size
	sort.Slice(allModels, func(i, j int) bool {
		return allModels[i].SizeBytes > allModels[j].SizeBytes
	})

	// Finalize
	sendEvent(models.ScanEvent{Type: models.EventCheckpoint, Phase: "complete", Message: "Compiling intelligence report...", Percent: 97})
	result := &models.ScanResult{
		Components: components,
		Models:     allModels,
		AgentCards: agentCards,
		MCPServers: mcpServers,
		ACPAgents:  acpAgents,
		Disks:      disks,
		GPU:        gpu,
		ScannedAt:  time.Now().Format(time.RFC3339),
	}

	// Known model dirs
	result.ModelDirs = getKnownModelDirs(home)
	s.cachedScan = result
	s.scanTime = time.Now()

	sendEvent(models.ScanEvent{
		Type:    models.EventComplete,
		Phase:   "complete",
		Message: "Scan complete.",
		Percent: 100,
		Detail:  result,
	})
}

// GPUHandler returns live GPU stats (for polling)
func (s *Server) GPUHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(scanner.ScanGPU())
}

// OSHandler returns OS info
func (s *Server) OSHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s.osInfo)
}

// CachedScanHandler returns last scan result if available
func (s *Server) CachedScanHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if s.cachedScan == nil {
		http.Error(w, `{"error":"no scan data yet"}`, http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(s.cachedScan)
}

// DeleteHandler handles model cleanup requests (SAFE MODE by default)
func (s *Server) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		FilePath string `json:"file_path"`
		Confirm  bool   `json:"confirm"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	log.Printf("[CLEANUP REQUEST] %s | Confirm: %v", req.FilePath, req.Confirm)
	log.Println("  -> SAFE MODE: File NOT deleted. Set AEYE_UNSAFE_DELETE=1 to enable.")

	// Uncomment to enable actual deletion:
	// if req.Confirm && os.Getenv("AEYE_UNSAFE_DELETE") == "1" {
	//     if err := os.Remove(req.FilePath); err != nil {
	//         http.Error(w, err.Error(), 500); return
	//     }
	// }

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "logged", "note": "safe mode active"})
}

func getKnownModelDirs(home string) []string {
	candidates := []string{
		home + "/.ollama/models",
		home + "/.cache/huggingface/hub",
		home + "/.cache/lm-studio/models",
		home + "/.local/share/lm-studio/models",
		home + "/Library/Application Support/LM Studio",
		home + "/AppData/Local/LMStudio",
	}
	var found []string
	for _, d := range candidates {
		if _, err := os.Stat(d); err == nil {
			found = append(found, d)
		}
	}
	return found
}

func (s *Server) Start(port string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/ws/scan", s.withCORS(s.ScanWebSocket))
	mux.HandleFunc("/api/gpu", s.withCORS(s.GPUHandler))
	mux.HandleFunc("/api/os", s.withCORS(s.OSHandler))
	mux.HandleFunc("/api/scan/cached", s.withCORS(s.CachedScanHandler))
	mux.HandleFunc("/api/delete", s.withCORS(s.DeleteHandler))
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status":"ok"}`))
	})

	fmt.Printf("\n  ░█▀█░█▀▀░█░█░█▀▀░░░█░░░█░█░█▀█░▀█▀░█▀▀░█░█░█▀▀░█▀▄\n")
	fmt.Printf("  ░█▀█░█▀▀░░█░░█▀▀░░░█▄▄░█▀█░█░█░░█░░█░░░█▀█░█▀▀░█▀▄\n")
	fmt.Printf("  ░▀░▀░▀▀▀░░▀░░▀▀▀░░░▀▀▀░▀░▀░▀░▀░░▀░░▀▀▀░▀░▀░▀▀▀░▀░▀\n\n")
	fmt.Printf("  Core API: http://localhost:%s\n", port)
	fmt.Printf("  OS:       %s %s (%s)\n", s.osInfo.Distro, s.osInfo.Version, s.osInfo.Arch)
	fmt.Printf("  GPU:      %v\n\n", s.osInfo.GPUAvailable)

	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Server error: %s", err)
	}
}
