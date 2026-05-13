package config

import (
	"strings"
	"time"
)

type ModelInfo struct {
	ID         string `json:"id"`
	Object     string `json:"object"`
	Created    int64  `json:"created"`
	OwnedBy    string `json:"owned_by"`
	Permission []any  `json:"permission,omitempty"`
}
type OllamaModelInfo struct {
	Name       string `json:"name"`
	Model      string `json:"model"`
	Size       int64  `json:"size"`
	ModifiedAt string `json:"modified_at"`
}
type OllamaCapabilitiesModelInfo struct {
	ID           string   `json:"id"`
	Capabilities []string `json:"capabilities"`
}

type ModelAliasReader interface {
	ModelAliases() map[string]string
}

const noThinkingModelSuffix = "-nothinking"

var deepSeekBaseModels = []ModelInfo{
	{ID: "deepseek-chat", Object: "model", Created: 1715635200, OwnedBy: "deepseek", Permission: []any{}},
	{ID: "deepseek-chat-search", Object: "model", Created: 1715635200, OwnedBy: "deepseek", Permission: []any{}},
	{ID: "deepseek-reasoner", Object: "model", Created: 1715635200, OwnedBy: "deepseek", Permission: []any{}},
	{ID: "deepseek-reasoner-search", Object: "model", Created: 1715635200, OwnedBy: "deepseek", Permission: []any{}},
	{ID: "deepseek-expert-chat", Object: "model", Created: 1715635200, OwnedBy: "deepseek", Permission: []any{}},
	{ID: "deepseek-expert-chat-search", Object: "model", Created: 1715635200, OwnedBy: "deepseek", Permission: []any{}},
	{ID: "deepseek-expert-reasoner", Object: "model", Created: 1715635200, OwnedBy: "deepseek", Permission: []any{}},
	{ID: "deepseek-expert-reasoner-search", Object: "model", Created: 1715635200, OwnedBy: "deepseek", Permission: []any{}},
	{ID: "deepseek-vision", Object: "model", Created: 1715635200, OwnedBy: "deepseek", Permission: []any{}},
	{ID: "deepseek-vision-search", Object: "model", Created: 1715635200, OwnedBy: "deepseek", Permission: []any{}},
	{ID: "deepseek-vision-reasoner", Object: "model", Created: 1715635200, OwnedBy: "deepseek", Permission: []any{}},
	{ID: "deepseek-vision-reasoner-search", Object: "model", Created: 1715635200, OwnedBy: "deepseek", Permission: []any{}},
}

var OllamaCapabilitiesModels = []OllamaCapabilitiesModelInfo{
	{ID: "deepseek-chat", Capabilities: []string{"tools"}},
	{ID: "deepseek-chat-search", Capabilities: []string{"tools"}},
	{ID: "deepseek-reasoner", Capabilities: []string{"tools", "thinking"}},
	{ID: "deepseek-reasoner-search", Capabilities: []string{"tools", "thinking"}},
	{ID: "deepseek-expert-chat", Capabilities: []string{"tools"}},
	{ID: "deepseek-expert-chat-search", Capabilities: []string{"tools"}},
	{ID: "deepseek-expert-reasoner", Capabilities: []string{"tools", "thinking"}},
	{ID: "deepseek-expert-reasoner-search", Capabilities: []string{"tools", "thinking"}},
	{ID: "deepseek-vision", Capabilities: []string{"tools", "vision"}},
	{ID: "deepseek-vision-search", Capabilities: []string{"tools", "vision"}},
	{ID: "deepseek-vision-reasoner", Capabilities: []string{"tools", "thinking", "vision"}},
	{ID: "deepseek-vision-reasoner-search", Capabilities: []string{"tools", "thinking", "vision"}},
}

var DeepSeekModels = deepSeekBaseModels
var OllamaModels = mapToOllamaModels(DeepSeekModels)
var claudeBaseModels = []ModelInfo{
	// Current aliases
	{ID: "claude-opus-4-6", Object: "model", Created: 1715635200, OwnedBy: "anthropic"},
	{ID: "claude-sonnet-4-6", Object: "model", Created: 1715635200, OwnedBy: "anthropic"},
	{ID: "claude-haiku-4-5", Object: "model", Created: 1715635200, OwnedBy: "anthropic"},

	// Claude 4.x snapshots and prior aliases kept for compatibility
	{ID: "claude-sonnet-4-5", Object: "model", Created: 1715635200, OwnedBy: "anthropic"},
	{ID: "claude-opus-4-1", Object: "model", Created: 1715635200, OwnedBy: "anthropic"},
	{ID: "claude-opus-4-1-20250805", Object: "model", Created: 1715635200, OwnedBy: "anthropic"},
	{ID: "claude-opus-4-0", Object: "model", Created: 1715635200, OwnedBy: "anthropic"},
	{ID: "claude-opus-4-20250514", Object: "model", Created: 1715635200, OwnedBy: "anthropic"},
	{ID: "claude-sonnet-4-5-20250929", Object: "model", Created: 1715635200, OwnedBy: "anthropic"},
	{ID: "claude-sonnet-4-0", Object: "model", Created: 1715635200, OwnedBy: "anthropic"},
	{ID: "claude-sonnet-4-20250514", Object: "model", Created: 1715635200, OwnedBy: "anthropic"},
	{ID: "claude-haiku-4-5-20251001", Object: "model", Created: 1715635200, OwnedBy: "anthropic"},

	// Claude 3.x (legacy/deprecated snapshots and aliases)
	{ID: "claude-3-7-sonnet-latest", Object: "model", Created: 1715635200, OwnedBy: "anthropic"},
	{ID: "claude-3-7-sonnet-20250219", Object: "model", Created: 1715635200, OwnedBy: "anthropic"},
	{ID: "claude-3-5-sonnet-latest", Object: "model", Created: 1715635200, OwnedBy: "anthropic"},
	{ID: "claude-3-5-sonnet-20240620", Object: "model", Created: 1715635200, OwnedBy: "anthropic"},
	{ID: "claude-3-5-sonnet-20241022", Object: "model", Created: 1715635200, OwnedBy: "anthropic"},
	{ID: "claude-3-opus-20240229", Object: "model", Created: 1715635200, OwnedBy: "anthropic"},
	{ID: "claude-3-sonnet-20240229", Object: "model", Created: 1715635200, OwnedBy: "anthropic"},
	{ID: "claude-3-5-haiku-latest", Object: "model", Created: 1715635200, OwnedBy: "anthropic"},
	{ID: "claude-3-5-haiku-20241022", Object: "model", Created: 1715635200, OwnedBy: "anthropic"},
	{ID: "claude-3-haiku-20240307", Object: "model", Created: 1715635200, OwnedBy: "anthropic"},
}

var ClaudeModels = appendNoThinkingVariants(claudeBaseModels)

func GetModelConfig(model string) (thinking bool, search bool, ok bool) {
	baseModel, noThinking := splitNoThinkingModel(model)
	if baseModel == "" {
		return false, false, false
	}
	switch baseModel {
	case "deepseek-chat", "deepseek-expert-chat", "deepseek-vision":
		return false, false, true
	case "deepseek-chat-search", "deepseek-expert-chat-search", "deepseek-vision-search":
		return false, true, true
	case "deepseek-reasoner", "deepseek-expert-reasoner", "deepseek-vision-reasoner":
		return !noThinking, false, true
	case "deepseek-reasoner-search", "deepseek-expert-reasoner-search", "deepseek-vision-reasoner-search":
		return !noThinking, true, true
	default:
		return false, false, false
	}
}

func GetModelType(model string) (modelType string, ok bool) {
	baseModel, _ := splitNoThinkingModel(model)
	switch baseModel {
	case "deepseek-chat", "deepseek-chat-search", "deepseek-reasoner", "deepseek-reasoner-search":
		return "default", true
	case "deepseek-expert-chat", "deepseek-expert-chat-search", "deepseek-expert-reasoner", "deepseek-expert-reasoner-search":
		return "expert", true
	case "deepseek-vision", "deepseek-vision-search", "deepseek-vision-reasoner", "deepseek-vision-reasoner-search":
		return "vision", true
	default:
		return "", false
	}
}

func IsSupportedDeepSeekModel(model string) bool {
	_, _, ok := GetModelConfig(model)
	return ok
}

func IsNoThinkingModel(model string) bool {
	_, noThinking := splitNoThinkingModel(model)
	return noThinking
}

func DefaultModelAliases() map[string]string {
	return map[string]string{
		// OpenAI GPT / ChatGPT families
		"chatgpt-4o":          "deepseek-chat",
		"gpt-4":               "deepseek-chat",
		"gpt-4-turbo":         "deepseek-chat",
		"gpt-4-turbo-preview": "deepseek-chat",
		"gpt-4.5-preview":     "deepseek-chat",
		"gpt-4o":              "deepseek-chat",
		"gpt-4o-mini":         "deepseek-chat",
		"gpt-4.1":             "deepseek-chat",
		"gpt-4.1-mini":        "deepseek-chat",
		"gpt-4.1-nano":        "deepseek-chat",
		"gpt-5":               "deepseek-chat",
		"gpt-5-chat":          "deepseek-chat",
		"gpt-5.1":             "deepseek-chat",
		"gpt-5.1-chat":        "deepseek-chat",
		"gpt-5.2":             "deepseek-chat",
		"gpt-5.2-chat":        "deepseek-chat",
		"gpt-5.3-chat":        "deepseek-chat",
		"gpt-5.4":             "deepseek-chat",
		"gpt-5.5":             "deepseek-chat",
		"gpt-5-mini":          "deepseek-chat",
		"gpt-5-nano":          "deepseek-chat",
		"gpt-5.4-mini":        "deepseek-chat",
		"gpt-5.4-nano":        "deepseek-chat",
		"gpt-5-pro":           "deepseek-expert-reasoner",
		"gpt-5.2-pro":         "deepseek-expert-reasoner",
		"gpt-5.4-pro":         "deepseek-expert-reasoner",
		"gpt-5.5-pro":         "deepseek-expert-reasoner",
		"gpt-5-codex":         "deepseek-expert-reasoner",
		"gpt-5.1-codex":       "deepseek-expert-reasoner",
		"gpt-5.1-codex-mini":  "deepseek-expert-reasoner",
		"gpt-5.1-codex-max":   "deepseek-expert-reasoner",
		"gpt-5.2-codex":       "deepseek-expert-reasoner",
		"gpt-5.3-codex":       "deepseek-expert-reasoner",
		"codex-mini-latest":   "deepseek-expert-reasoner",

		// OpenAI reasoning / research families
		"o1":                    "deepseek-expert-reasoner",
		"o1-preview":            "deepseek-expert-reasoner",
		"o1-mini":               "deepseek-expert-reasoner",
		"o1-pro":                "deepseek-expert-reasoner",
		"o3":                    "deepseek-expert-reasoner",
		"o3-mini":               "deepseek-expert-reasoner",
		"o3-pro":                "deepseek-expert-reasoner",
		"o3-deep-research":      "deepseek-expert-reasoner-search",
		"o4-mini":               "deepseek-expert-reasoner",
		"o4-mini-deep-research": "deepseek-expert-reasoner-search",

		// Claude current and historical aliases
		"claude-opus-4-6":            "deepseek-expert-reasoner",
		"claude-opus-4-1":            "deepseek-expert-reasoner",
		"claude-opus-4-1-20250805":   "deepseek-expert-reasoner",
		"claude-opus-4-0":            "deepseek-expert-reasoner",
		"claude-opus-4-20250514":     "deepseek-expert-reasoner",
		"claude-sonnet-4-6":          "deepseek-chat",
		"claude-sonnet-4-5":          "deepseek-chat",
		"claude-sonnet-4-5-20250929": "deepseek-chat",
		"claude-sonnet-4-0":          "deepseek-chat",
		"claude-sonnet-4-20250514":   "deepseek-chat",
		"claude-haiku-4-5":           "deepseek-chat",
		"claude-haiku-4-5-20251001":  "deepseek-chat",
		"claude-3-7-sonnet":          "deepseek-chat",
		"claude-3-7-sonnet-latest":   "deepseek-chat",
		"claude-3-7-sonnet-20250219": "deepseek-chat",
		"claude-3-5-sonnet":          "deepseek-chat",
		"claude-3-5-sonnet-latest":   "deepseek-chat",
		"claude-3-5-sonnet-20240620": "deepseek-chat",
		"claude-3-5-sonnet-20241022": "deepseek-chat",
		"claude-3-5-haiku":           "deepseek-chat",
		"claude-3-5-haiku-latest":    "deepseek-chat",
		"claude-3-5-haiku-20241022":  "deepseek-chat",
		"claude-3-opus":              "deepseek-expert-reasoner",
		"claude-3-opus-20240229":     "deepseek-expert-reasoner",
		"claude-3-sonnet":            "deepseek-chat",
		"claude-3-sonnet-20240229":   "deepseek-chat",
		"claude-3-haiku":             "deepseek-chat",
		"claude-3-haiku-20240307":    "deepseek-chat",

		// Gemini current and historical text / multimodal models
		"gemini-pro":            "deepseek-expert-reasoner",
		"gemini-pro-vision":     "deepseek-vision",
		"gemini-pro-latest":     "deepseek-expert-reasoner",
		"gemini-flash-latest":   "deepseek-chat",
		"gemini-1.5-pro":        "deepseek-expert-reasoner",
		"gemini-1.5-flash":      "deepseek-chat",
		"gemini-1.5-flash-8b":   "deepseek-chat",
		"gemini-2.0-flash":      "deepseek-chat",
		"gemini-2.0-flash-lite": "deepseek-chat",
		"gemini-2.5-pro":        "deepseek-expert-reasoner",
		"gemini-2.5-flash":      "deepseek-chat",
		"gemini-2.5-flash-lite": "deepseek-chat",
		"gemini-3.1-pro":        "deepseek-expert-reasoner",
		"gemini-3-pro":          "deepseek-expert-reasoner",
		"gemini-3-flash":        "deepseek-chat",
		"gemini-3.1-flash":      "deepseek-chat",
		"gemini-3.1-flash-lite": "deepseek-chat",

		"llama-3.1-70b-instruct": "deepseek-chat",
		"qwen-max":               "deepseek-chat",
	}
}

func ResolveModel(store ModelAliasReader, requested string) (string, bool) {
	model := lower(strings.TrimSpace(requested))
	if model == "" {
		return "", false
	}
	aliases := loadModelAliases(store)
	if IsSupportedDeepSeekModel(model) {
		return model, true
	}
	if mapped, ok := aliases[model]; ok && IsSupportedDeepSeekModel(mapped) {
		return mapped, true
	}
	baseModel, noThinking := splitNoThinkingModel(model)
	if mapped, ok := aliases[baseModel]; ok && IsSupportedDeepSeekModel(mapped) {
		return withNoThinkingVariant(mapped, noThinking), true
	}
	return "", false
}

func lower(s string) string {
	b := []byte(s)
	for i, c := range b {
		if c >= 'A' && c <= 'Z' {
			b[i] = c + 32
		}
	}
	return string(b)
}

func OpenAIModelsResponse() map[string]any {
	return map[string]any{"object": "list", "data": DeepSeekModels}
}

func OpenAIModelByID(store ModelAliasReader, id string) (ModelInfo, bool) {
	canonical, ok := ResolveModel(store, id)
	if !ok {
		return ModelInfo{}, false
	}
	for _, model := range DeepSeekModels {
		if model.ID == canonical {
			return model, true
		}
	}
	return ModelInfo{}, false
}

func OllamaModelsResponse() map[string]any {
	return map[string]any{"models": OllamaModels}
}

func OllamaModelByID(store ModelAliasReader, id string) (OllamaCapabilitiesModelInfo, bool) {
	canonical, ok := ResolveModel(store, id)
	if !ok {
		return OllamaCapabilitiesModelInfo{}, false
	}
	for _, model := range OllamaCapabilitiesModels {
		if model.ID == canonical {
			return model, true
		}
	}
	return OllamaCapabilitiesModelInfo{}, false
}

func ClaudeModelsResponse() map[string]any {
	resp := map[string]any{"object": "list", "data": ClaudeModels}
	if len(ClaudeModels) > 0 {
		resp["first_id"] = ClaudeModels[0].ID
		resp["last_id"] = ClaudeModels[len(ClaudeModels)-1].ID
	} else {
		resp["first_id"] = nil
		resp["last_id"] = nil
	}
	resp["has_more"] = false
	return resp
}

func appendNoThinkingVariants(models []ModelInfo) []ModelInfo {
	out := make([]ModelInfo, 0, len(models)*2)
	for _, model := range models {
		out = append(out, model)
		variant := model
		variant.ID = withNoThinkingVariant(model.ID, true)
		out = append(out, variant)
	}
	return out
}
func mapToOllamaModels(models []ModelInfo) []OllamaModelInfo {
	out := make([]OllamaModelInfo, 0, len(models))
	for _, model := range models {
		var modifiedAt string
		if model.Created > 0 {
			modifiedAt = time.Unix(model.Created, 0).Format(time.RFC3339)
		}
		ollamaModel := OllamaModelInfo{
			Name:       model.ID,
			Model:      model.ID,
			Size:       0,
			ModifiedAt: modifiedAt,
		}
		out = append(out, ollamaModel)
	}
	return out
}

func splitNoThinkingModel(model string) (string, bool) {
	model = lower(strings.TrimSpace(model))
	if strings.HasSuffix(model, noThinkingModelSuffix) {
		return strings.TrimSuffix(model, noThinkingModelSuffix), true
	}
	return model, false
}

func withNoThinkingVariant(model string, enabled bool) string {
	baseModel, _ := splitNoThinkingModel(model)
	if !enabled {
		return baseModel
	}
	if baseModel == "" {
		return ""
	}
	return baseModel + noThinkingModelSuffix
}

func loadModelAliases(store ModelAliasReader) map[string]string {
	aliases := DefaultModelAliases()
	if store != nil {
		for k, v := range store.ModelAliases() {
			aliases[lower(strings.TrimSpace(k))] = lower(strings.TrimSpace(v))
		}
	}
	return aliases
}
