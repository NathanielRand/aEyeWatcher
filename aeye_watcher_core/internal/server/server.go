package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
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
	s := &Server{osInfo: os_detect.Detect()}
	s.loadCachedScan() // Load persisted cache on startup
	return s
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

// ScanWebSocket streams live scan events to the client.
func (s *Server) ScanWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WS upgrade error:", err)
		return
	}
	defer conn.Close()

	scanStart := time.Now()
	sendEvent := func(evt models.ScanEvent) {
		data, _ := json.Marshal(evt)
		conn.WriteMessage(websocket.TextMessage, data)
		PrintScanEvent(evt, time.Since(scanStart))
	}

	home, _ := os.UserHomeDir()

	// Phase 1: OS & System
	sendEvent(models.ScanEvent{Type: models.EventCheckpoint, Phase: "system", Message: "Initializing system profile...", Percent: 2})
	time.Sleep(150 * time.Millisecond)
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
		if basePercent < 44 {
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
	sendEvent(models.ScanEvent{Type: models.EventCheckpoint, Phase: "mcp", Message: "Hunting MCP server configurations...", Percent: 47})
	mcpCh := make(chan string, 20)
	var mcpServers []models.MCPServer
	go func() {
		mcpServers = scanner.ScanMCPServers(home, mcpCh)
		close(mcpCh)
	}()
	for msg := range mcpCh {
		sendEvent(models.ScanEvent{Type: models.EventProgress, Phase: "mcp", Message: msg, Percent: 50})
	}
	sendEvent(models.ScanEvent{Type: models.EventFound, Phase: "mcp", Message: fmt.Sprintf("Found %d MCP server configuration(s)", len(mcpServers)), Percent: 54})

	// Phase 6: ACP Agents
	sendEvent(models.ScanEvent{Type: models.EventCheckpoint, Phase: "acp", Message: "Probing ACP agent endpoints...", Percent: 56})
	acpAgents := scanner.ScanACPAgents(nil)
	sendEvent(models.ScanEvent{Type: models.EventFound, Phase: "acp", Message: fmt.Sprintf("Found %d ACP agent(s)", len(acpAgents)), Percent: 59})

	// Phase 7: Models
	sendEvent(models.ScanEvent{Type: models.EventCheckpoint, Phase: "models", Message: "Parsing Ollama model manifests...", Percent: 61})
	ms := scanner.NewModelScanner(home, nil)
	ollamaModels := ms.ScanOllamaManifests()
	sendEvent(models.ScanEvent{Type: models.EventFound, Phase: "models", Message: fmt.Sprintf("Ollama: %d model(s)", len(ollamaModels)), Percent: 64})

	sendEvent(models.ScanEvent{Type: models.EventCheckpoint, Phase: "models", Message: "Scanning LM Studio hub...", Percent: 65})
	lmStudioModels := ms.ScanLMStudioModels()
	sendEvent(models.ScanEvent{Type: models.EventFound, Phase: "models", Message: fmt.Sprintf("LM Studio: %d model(s)", len(lmStudioModels)), Percent: 67})

	sendEvent(models.ScanEvent{Type: models.EventCheckpoint, Phase: "models", Message: "Parsing HuggingFace cache...", Percent: 68})
	hfModels := ms.ScanHuggingFaceCache()
	sendEvent(models.ScanEvent{Type: models.EventFound, Phase: "models", Message: fmt.Sprintf("HuggingFace: %d model(s)", len(hfModels)), Percent: 70})

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
		if msgCount%5 == 0 {
			sendEvent(models.ScanEvent{Type: models.EventProgress, Phase: "models", Message: msg, Percent: deepPercent})
			if deepPercent < 84 {
				deepPercent++
			}
		}
	}

	allModels := scanner.DeduplicateModels(append(append(append(ollamaModels, lmStudioModels...), hfModels...), deepModels...))
	orphaned := 0
	for _, m := range allModels {
		if m.IsOrphaned {
			orphaned++
		}
	}
	sendEvent(models.ScanEvent{Type: models.EventFound, Phase: "models", Message: fmt.Sprintf("Total models: %d (%d orphaned)", len(allModels), orphaned), Percent: 85})

	// Phase 8: Agent Cards
	sendEvent(models.ScanEvent{Type: models.EventCheckpoint, Phase: "agents", Message: "Building agent ecosystem profiles...", Percent: 87})
	agentCh := make(chan string, 30)
	var agentCards []models.AgentCard
	go func() {
		agentCards = scanner.BuildAgentCards(home, agentCh)
		close(agentCh)
	}()
	for msg := range agentCh {
		sendEvent(models.ScanEvent{Type: models.EventProgress, Phase: "agents", Message: msg, Percent: 90})
	}
	sendEvent(models.ScanEvent{Type: models.EventFound, Phase: "agents", Message: fmt.Sprintf("Profiled %d agent ecosystems", len(agentCards)), Percent: 93})

	sort.Slice(allModels, func(i, j int) bool {
		return allModels[i].SizeBytes > allModels[j].SizeBytes
	})

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
	result.ModelDirs = getKnownModelDirs(home)
	s.cachedScan = result
	s.scanTime = time.Now()

	// Persist the results
	s.saveCachedScan(result)
	s.addToHistory(result)

	sendEvent(models.ScanEvent{
		Type:    models.EventComplete,
		Phase:   "complete",
		Message: "Scan complete.",
		Percent: 100,
		Detail:  result,
	})
}

func (s *Server) GPUHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(scanner.ScanGPU())
}

func (s *Server) OSHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s.osInfo)
}

// CACHING LOGIC
const scanCacheFile = "scan_cache.json"

func (s *Server) loadCachedScan() {
	data, err := os.ReadFile(scanCacheFile)
	if err != nil {
		return // No cache yet
	}
	var cached models.ScanResult
	if err := json.Unmarshal(data, &cached); err == nil {
		s.cachedScan = &cached
		log.Printf("Loaded cached scan from %s", scanCacheFile)
	}
}

func (s *Server) saveCachedScan(result *models.ScanResult) {
	data, _ := json.MarshalIndent(result, "", "  ")
	os.WriteFile(scanCacheFile, data, 0644)
}

func (s *Server) CachedScanHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if s.cachedScan == nil {
		http.Error(w, `{"error":"no scan data yet"}`, http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(s.cachedScan)
}

// SCAN HISTORY LOGIC
const scanHistoryDir = "scan_history"

func (s *Server) addToHistory(result *models.ScanResult) {
	os.MkdirAll(scanHistoryDir, 0755)
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	path := filepath.Join(scanHistoryDir, fmt.Sprintf("scan_%s.json", timestamp))
	data, _ := json.MarshalIndent(result, "", "  ")
	os.WriteFile(path, data, 0644)

	// Keep last 50 scans
	cleanupOldScans(50)
}

func cleanupOldScans(keep int) {
	entries, _ := os.ReadDir(scanHistoryDir)
	if len(entries) <= keep {
		return
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})
	for i := 0; i < len(entries)-keep; i++ {
		os.Remove(filepath.Join(scanHistoryDir, entries[i].Name()))
	}
}

func (s *Server) ScanHistoryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	entries, _ := os.ReadDir(scanHistoryDir)
	var history []struct {
		File      string `json:"file"`
		Timestamp string `json:"timestamp"`
		Size      int64  `json:"size"`
	}
	for _, e := range entries {
		info, _ := e.Info()
		history = append(history, struct {
			File      string `json:"file"`
			Timestamp string `json:"timestamp"`
			Size      int64  `json:"size"`
		}{
			File:      e.Name(),
			Timestamp: info.ModTime().Format(time.RFC3339),
			Size:      info.Size(),
		})
	}
	// Reverse so newest first
	for i, j := 0, len(history)-1; i < j; i, j = j, i {
		history[i], history[j] = history[j], history[i]
	}
	json.NewEncoder(w).Encode(history)
}

func (s *Server) ScanHistoryItemHandler(w http.ResponseWriter, r *http.Request) {
	// GET /api/scan/history/{filename}
	file := strings.TrimPrefix(r.URL.Path, "/api/scan/history/")
	data, err := os.ReadFile(filepath.Join(scanHistoryDir, file))
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

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
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "logged", "note": "safe mode active"})
}

// OpenPathHandler opens a path in the OS default file manager / app.
func (s *Server) OpenPathHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST only", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Path string `json:"path"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Path == "" {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	go func() {
		switch runtime.GOOS {
		case "windows":
			exec.Command("explorer", req.Path).Start()
		case "darwin":
			exec.Command("open", req.Path).Start()
		default:
			exec.Command("xdg-open", req.Path).Start()
		}
	}()
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"ok"}`))
}

// ReadFileHandler returns the text content of a readable file (<1MB).
func (s *Server) ReadFileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST only", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Path string `json:"path"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Path == "" {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	ext := strings.ToLower(filepath.Ext(req.Path))
	allowed := map[string]bool{
		".md": true, ".json": true, ".yaml": true, ".yml": true,
		".txt": true, ".toml": true, ".env": true, ".conf": true,
		".cfg": true, ".ini": true, ".log": true, ".soul": true,
	}
	if !allowed[ext] {
		http.Error(w, "file type not readable", http.StatusForbidden)
		return
	}
	info, err := os.Stat(req.Path)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	if info.Size() > 1_048_576 {
		http.Error(w, "file too large", http.StatusRequestEntityTooLarge)
		return
	}
	data, err := os.ReadFile(req.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write(data)
}

func getKnownModelDirs(home string) []string {
	candidates := []string{
		home + "/.ollama/models",
		home + "/.lmstudio/hub/models",
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
	mux.HandleFunc("/api/scan/history", s.withCORS(s.ScanHistoryHandler))
	mux.HandleFunc("/api/scan/history/", s.withCORS(s.ScanHistoryItemHandler))
	mux.HandleFunc("/api/delete", s.withCORS(s.DeleteHandler))
	mux.HandleFunc("/api/open", s.withCORS(s.OpenPathHandler))
	mux.HandleFunc("/api/read", s.withCORS(s.ReadFileHandler))
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status":"ok"}`))
	})

	PrintBanner(s.osInfo, port)

	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Server error: %s", err)
	}
}
