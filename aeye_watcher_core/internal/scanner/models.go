package scanner

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/aeye/watcher-core/internal/models"
)

// Known AI model file extensions and their formats
var modelExtensions = map[string]string{
	".gguf":        "GGUF",
	".ggml":        "GGML",
	".bin":         "PyTorch/GGML",
	".safetensors": "Safetensors",
	".onnx":        "ONNX",
	".pt":          "PyTorch",
	".pth":         "PyTorch",
	".pkl":         "Pickle/PyTorch",
	".mlmodel":     "CoreML",
	".tflite":      "TensorFlow Lite",
	".pb":          "TensorFlow",
	".engine":      "TensorRT",
	".llamafile":   "Llamafile",
}

// Known provider model directories to mark as non-orphaned
var knownProviderDirs = []string{
	".ollama",
	".cache/huggingface",
	".cache/lm-studio",
	".local/share/lm-studio",
	"Library/Application Support/LM Studio", // macOS
	"AppData/Local/LMStudio",                // Windows
	"AppData/Roaming/LMStudio",
	".cache/torch",
	".cache/whisper",
	"koboldcpp",
	"llama.cpp",
	"text-generation-webui/models",
	"oobabooga",
}

// Directories to always skip during deep scan
var skipDirs = map[string]bool{
	"node_modules": true,
	".npm":         true,
	".yarn":        true,
	"go":           true,
	".git":         true,
	".vscode":      true,
	"__pycache__":  true,
	".Trash":       true,
	"$Recycle.Bin": true,
	"Windows":      true,
	"System32":     true,
	"SysWOW64":     true,
}

// Minimum file size for model files (avoid tiny test files) - 10MB
const minModelSize = 10 * 1024 * 1024

type ModelScanner struct {
	home       string
	knownDirs  []string
	progressCh chan<- string
}

func NewModelScanner(home string, progressCh chan<- string) *ModelScanner {
	knownDirs := make([]string, len(knownProviderDirs))
	for i, d := range knownProviderDirs {
		knownDirs[i] = filepath.Join(home, d)
	}
	return &ModelScanner{home: home, knownDirs: knownDirs, progressCh: progressCh}
}

func (ms *ModelScanner) isKnownProviderPath(path string) bool {
	normalized := filepath.ToSlash(path)
	for _, known := range knownProviderDirs {
		if strings.Contains(normalized, known) {
			return true
		}
	}
	return false
}

// ScanOllamaManifests parses the Ollama manifest format for accurate model info
func (ms *ModelScanner) ScanOllamaManifests() []models.AIModel {
	var result []models.AIModel
	registryBase := filepath.Join(ms.home, ".ollama", "models", "manifests", "registry.ollama.ai")
	namespaces, err := os.ReadDir(registryBase)
	if err != nil {
		return result
	}

	for _, ns := range namespaces {
		if !ns.IsDir() {
			continue
		}
		nsPath := filepath.Join(registryBase, ns.Name())
		modelEntries, err := os.ReadDir(nsPath)
		if err != nil {
			continue
		}

		for _, modelEntry := range modelEntries {
			if !modelEntry.IsDir() {
				continue
			}
			modelName := modelEntry.Name()
			if ns.Name() != "library" {
				modelName = ns.Name() + "/" + modelName
			}

			modelPath := filepath.Join(nsPath, modelEntry.Name())
			tagEntries, err := os.ReadDir(modelPath)
			if err != nil {
				continue
			}

			for _, tagEntry := range tagEntries {
				if tagEntry.IsDir() {
					continue
				}
				manifestPath := filepath.Join(modelPath, tagEntry.Name())
				data, err := os.ReadFile(manifestPath)
				if err != nil {
					continue
				}
				var manifest struct {
					Layers []struct {
						Size int64 `json:"size"`
					} `json:"layers"`
				}
				if json.Unmarshal(data, &manifest) == nil {
					var totalSize int64
					for _, layer := range manifest.Layers {
						totalSize += layer.Size
					}
					result = append(result, models.AIModel{
						Name: modelName, Tag: tagEntry.Name(), SizeBytes: totalSize,
						Provider: "Ollama", FilePath: manifestPath,
						Extension: ".manifest", Format: "Ollama Manifest", IsOrphaned: false,
					})
				}
			}
		}
	}
	return result
}

// ScanHuggingFaceCache parses the HuggingFace cache directory structure
func (ms *ModelScanner) ScanHuggingFaceCache() []models.AIModel {
	var result []models.AIModel
	hfCache := filepath.Join(ms.home, ".cache", "huggingface", "hub")
	if _, err := os.Stat(hfCache); os.IsNotExist(err) {
		return result
	}

	entries, err := os.ReadDir(hfCache)
	if err != nil {
		return result
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		// HF cache structure: models--org--name / snapshots / hash / files
		if !strings.HasPrefix(entry.Name(), "models--") {
			continue
		}
		parts := strings.SplitN(strings.TrimPrefix(entry.Name(), "models--"), "--", 2)
		modelName := entry.Name()
		if len(parts) == 2 {
			modelName = parts[0] + "/" + parts[1]
		}

		var totalSize int64
		filepath.WalkDir(filepath.Join(hfCache, entry.Name()), func(path string, d os.DirEntry, err error) error {
			if err == nil && !d.IsDir() {
				info, err := d.Info()
				if err == nil {
					totalSize += info.Size()
				}
			}
			return nil
		})

		if totalSize > minModelSize {
			result = append(result, models.AIModel{
				Name: modelName, Tag: "latest", SizeBytes: totalSize,
				Provider: "HuggingFace", FilePath: filepath.Join(hfCache, entry.Name()),
				Format: "HuggingFace Cache", IsOrphaned: false,
			})
		}
	}
	return result
}

// DeepScan walks the filesystem looking for AI model files outside known dirs
func (ms *ModelScanner) DeepScan(progressCh chan<- string) []models.AIModel {
	var result []models.AIModel

	// Determine scan roots based on OS
	roots := getScanRoots(ms.home)

	for _, root := range roots {
		if _, err := os.Stat(root); os.IsNotExist(err) {
			continue
		}
		if progressCh != nil {
			progressCh <- "Scanning " + root
		}

		filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return nil
			}

			if d.IsDir() {
				name := d.Name()
				if skipDirs[name] {
					return filepath.SkipDir
				}
				// Skip very deep paths in certain dirs
				rel, _ := filepath.Rel(ms.home, path)
				depth := len(strings.Split(rel, string(os.PathSeparator)))
				if depth > 12 {
					return filepath.SkipDir
				}
				return nil
			}

			ext := strings.ToLower(filepath.Ext(d.Name()))
			format, ok := modelExtensions[ext]
			if !ok {
				return nil
			}

			info, err := d.Info()
			if err != nil || info.Size() < minModelSize {
				return nil
			}

			isOrphaned := !ms.isKnownProviderPath(path)
			name := strings.TrimSuffix(d.Name(), filepath.Ext(d.Name()))

			// Determine provider from path
			provider := guessProvider(path)

			result = append(result, models.AIModel{
				Name: name, Tag: "local", SizeBytes: info.Size(),
				Provider: provider, FilePath: path,
				Extension: ext, Format: format, IsOrphaned: isOrphaned,
			})
			return nil
		})
	}

	return result
}

func getScanRoots(home string) []string {
	switch runtime.GOOS {
	case "windows":
		return []string{
			home,
			"C:\\ProgramData\\ollama",
			"C:\\ProgramData\\LMStudio",
		}
	case "darwin":
		return []string{
			home,
			"/Applications",
		}
	default: // linux
		return []string{
			home,
			"/opt",
			"/usr/local/share",
		}
	}
}

func guessProvider(path string) string {
	p := filepath.ToSlash(strings.ToLower(path))
	switch {
	case strings.Contains(p, "ollama"):
		return "Ollama"
	case strings.Contains(p, "huggingface") || strings.Contains(p, "hf_cache"):
		return "HuggingFace"
	case strings.Contains(p, "lm-studio") || strings.Contains(p, "lmstudio"):
		return "LM Studio"
	case strings.Contains(p, "kobold"):
		return "KoboldCPP"
	case strings.Contains(p, "llama.cpp") || strings.Contains(p, "llamacpp"):
		return "llama.cpp"
	case strings.Contains(p, "text-generation-webui") || strings.Contains(p, "oobabooga"):
		return "Text Gen WebUI"
	case strings.Contains(p, "whisper"):
		return "Whisper"
	case strings.Contains(p, "torch") || strings.Contains(p, "pytorch"):
		return "PyTorch Cache"
	default:
		return "Unknown / Orphaned"
	}
}

// DeduplicateModels removes duplicate entries (same file path)
func DeduplicateModels(dupedModels []models.AIModel) []models.AIModel {
	seen := map[string]bool{}
	var unique []models.AIModel
	for _, m := range dupedModels {
		if !seen[m.FilePath] {
			seen[m.FilePath] = true
			unique = append(unique, m)
		}
	}
	return unique
}
