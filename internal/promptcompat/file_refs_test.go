package promptcompat

import (
	"strings"
	"testing"
)

func TestCollectOpenAIRefFileIDsFromFileTags(t *testing.T) {
	req := map[string]any{
		"model": "deepseek-chat",
		"messages": []any{
			map[string]any{
				"role":    "user",
				"content": "reuse <||file:report.txt:file-abc:user@example.com||> please",
			},
		},
	}

	ids := CollectOpenAIRefFileIDs(req)
	if len(ids) != 1 || ids[0] != "file-abc" {
		t.Fatalf("unexpected file ids: %#v", ids)
	}

	accounts := CollectOpenAIFileTagAccounts(req)
	if len(accounts) != 1 || accounts[0] != "user@example.com" {
		t.Fatalf("unexpected accounts: %#v", accounts)
	}
}

func TestStripOpenAIFileTags(t *testing.T) {
	req := map[string]any{
		"messages": []any{
			map[string]any{
				"role":    "user",
				"content": "before <||file:report.txt:file-abc:user@example.com||> after",
			},
		},
	}

	StripOpenAIFileTags(req)
	messages := req["messages"].([]any)
	msg := messages[0].(map[string]any)
	content := msg["content"].(string)
	if strings.Contains(content, "<||file:") || strings.Contains(content, "file-abc") {
		t.Fatalf("expected file tag stripped, got %q", content)
	}
	if content != "before  after" {
		t.Fatalf("unexpected stripped content: %q", content)
	}
}

func TestNormalizeOpenAIChatRequestCollectsAndStripsFileTags(t *testing.T) {
	req := map[string]any{
		"model": "deepseek-chat",
		"messages": []any{
			map[string]any{
				"role":    "user",
				"content": "read <||file:report.txt:file-abc:13800138000||> now",
			},
		},
	}

	out, err := NormalizeOpenAIChatRequest(mockPromptConfig{}, req, "")
	if err != nil {
		t.Fatalf("NormalizeOpenAIChatRequest error: %v", err)
	}
	if len(out.RefFileIDs) != 1 || out.RefFileIDs[0] != "file-abc" {
		t.Fatalf("unexpected ref file ids: %#v", out.RefFileIDs)
	}
	if strings.Contains(out.FinalPrompt, "<||file:") || strings.Contains(out.FinalPrompt, "file-abc") {
		t.Fatalf("file tag leaked into prompt: %q", out.FinalPrompt)
	}
}

type mockPromptConfig struct{}

func (mockPromptConfig) ModelAliases() map[string]string { return nil }
