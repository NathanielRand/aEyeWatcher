package models

// ---- Scan Progress Events ----

type ScanEventType string

const (
	EventCheckpoint ScanEventType = "checkpoint"
	EventProgress   ScanEventType = "progress"
	EventFound      ScanEventType = "found"
	EventComplete   ScanEventType = "complete"
	EventError      ScanEventType = "error"
)

type ScanEvent struct {
	Type    ScanEventType `json:"type"`
	Phase   string        `json:"phase"`
	Message string        `json:"message"`
	Percent int           `json:"percent"`
	Detail  interface{}   `json:"detail,omitempty"`
}

// ---- Components ----

type ComponentStatus string

const (
	StatusRunning   ComponentStatus = "Running"
	StatusInstalled ComponentStatus = "Installed"
	StatusNotFound  ComponentStatus = "Not Found"
	StatusDegraded  ComponentStatus = "Degraded"
)

type Component struct {
	Name        string          `json:"name"`
	Type        string          `json:"type"`        // Provider, Agent, Gateway, Vector DB, GUI Client, TTS/Audio, MCP Server, ACP Agent
	Status      ComponentStatus `json:"status"`
	Path        string          `json:"path,omitempty"`
	Port        string          `json:"port,omitempty"`
	Version     string          `json:"version,omitempty"`
	Description string          `json:"description,omitempty"`
	URL         string          `json:"url,omitempty"`
	LastSeen    string          `json:"last_seen,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// ---- AI Models ----

type AIModel struct {
	Name       string `json:"name"`
	Tag        string `json:"tag"`
	SizeBytes  int64  `json:"size_bytes"`
	Provider   string `json:"provider"`
	FilePath   string `json:"file_path"`
	Extension  string `json:"extension"`
	IsOrphaned bool   `json:"is_orphaned"` // found outside known provider dirs
	Format     string `json:"format"`      // GGUF, Safetensors, ONNX, PyTorch, etc.
}

// ---- Agent Ecosystem ----

type AgentFile struct {
	Type     string `json:"type"`     // skill, soul, config, memory, tool, persona
	Name     string `json:"name"`
	Path     string `json:"path"`
	SizeBytes int64 `json:"size_bytes"`
}

type AgentCard struct {
	Name        string      `json:"name"`
	Ecosystem   string      `json:"ecosystem"` // claude, openai, openinterpreter, etc.
	Status      string      `json:"status"`    // detected, partial, not_found
	Description string      `json:"description"`
	Icon        string      `json:"icon"`
	Files       []AgentFile `json:"files"`
	ExpectedFiles []ExpectedFile `json:"expected_files"`
}

type ExpectedFile struct {
	Label   string `json:"label"`
	Pattern string `json:"pattern"`
	Found   bool   `json:"found"`
	Path    string `json:"path,omitempty"`
}

// ---- MCP / ACP ----

type MCPServer struct {
	Name     string `json:"name"`
	URL      string `json:"url,omitempty"`
	Config   string `json:"config,omitempty"` // path to config file
	Status   string `json:"status"`
	Transport string `json:"transport"` // stdio, http, sse
	Tools    []string `json:"tools,omitempty"`
}

type ACPAgent struct {
	Name     string `json:"name"`
	Endpoint string `json:"endpoint"`
	Status   string `json:"status"`
	Protocol string `json:"protocol"` // ACP v0.x
}

// ---- Disk ----

type DiskInfo struct {
	Mount     string  `json:"mount"`
	Device    string  `json:"device"`
	FSType    string  `json:"fs_type"`
	TotalGB   float64 `json:"total_gb"`
	UsedGB    float64 `json:"used_gb"`
	FreeGB    float64 `json:"free_gb"`
	UsePct    float64 `json:"use_pct"`
}

// ---- GPU ----

type GPUInfo struct {
	Name      string  `json:"name,omitempty"`
	UsedMB    int     `json:"used_mb"`
	TotalMB   int     `json:"total_mb"`
	UsePct    float64 `json:"use_pct"`
	TempC     int     `json:"temp_c,omitempty"`
	Error     string  `json:"error,omitempty"`
	Backend   string  `json:"backend"` // nvidia, amd, none
}

// ---- Full Scan Result ----

type ScanResult struct {
	Components []Component `json:"components"`
	Models     []AIModel   `json:"models"`
	AgentCards []AgentCard `json:"agent_cards"`
	MCPServers []MCPServer `json:"mcp_servers"`
	ACPAgents  []ACPAgent  `json:"acp_agents"`
	ModelDirs  []string    `json:"model_dirs"`
	Disks      []DiskInfo  `json:"disks"`
	GPU        GPUInfo     `json:"gpu"`
	ScannedAt  string      `json:"scanned_at"`
}
