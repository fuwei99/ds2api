package config

import "testing"

type mockModelAliasReader map[string]string

func (m mockModelAliasReader) ModelAliases() map[string]string { return m }

func TestResolveModelDirectDeepSeek(t *testing.T) {
	got, ok := ResolveModel(nil, "deepseek-chat")
	if !ok || got != "deepseek-chat" {
		t.Fatalf("expected deepseek-chat, got ok=%v model=%q", ok, got)
	}
}

func TestResolveModelDirectDeepSeekSearch(t *testing.T) {
	got, ok := ResolveModel(nil, "deepseek-chat-search")
	if !ok || got != "deepseek-chat-search" {
		t.Fatalf("expected deepseek-chat-search, got ok=%v model=%q", ok, got)
	}
}

func TestResolveModelAlias(t *testing.T) {
	got, ok := ResolveModel(nil, "gpt-4.1")
	if !ok || got != "deepseek-chat" {
		t.Fatalf("expected alias gpt-4.1 -> deepseek-chat, got ok=%v model=%q", ok, got)
	}
}

func TestResolveLatestOpenAIAlias(t *testing.T) {
	got, ok := ResolveModel(nil, "gpt-5.5")
	if !ok || got != "deepseek-chat" {
		t.Fatalf("expected alias gpt-5.5 -> deepseek-chat, got ok=%v model=%q", ok, got)
	}
}

func TestResolveLatestClaudeAlias(t *testing.T) {
	got, ok := ResolveModel(nil, "claude-sonnet-4-6")
	if !ok || got != "deepseek-chat" {
		t.Fatalf("expected alias claude-sonnet-4-6 -> deepseek-chat, got ok=%v model=%q", ok, got)
	}
}

func TestResolveLatestClaudeAliasNoThinking(t *testing.T) {
	got, ok := ResolveModel(nil, "claude-sonnet-4-6-nothinking")
	if !ok || got != "deepseek-chat-nothinking" {
		t.Fatalf("expected alias claude-sonnet-4-6-nothinking -> deepseek-chat-nothinking, got ok=%v model=%q", ok, got)
	}
}

func TestResolveExpandedHistoricalAliases(t *testing.T) {
	cases := []struct {
		name  string
		model string
		want  string
	}{
		{name: "openai old chatgpt", model: "chatgpt-4o", want: "deepseek-chat"},
		{name: "openai codex max", model: "gpt-5.1-codex-max", want: "deepseek-expert-reasoner"},
		{name: "openai deep research", model: "o3-deep-research", want: "deepseek-expert-reasoner-search"},
		{name: "openai historical reasoning", model: "o1-preview", want: "deepseek-expert-reasoner"},
		{name: "claude latest historical", model: "claude-3-5-sonnet-latest", want: "deepseek-chat"},
		{name: "claude historical opus", model: "claude-3-opus-20240229", want: "deepseek-expert-reasoner"},
		{name: "claude historical haiku", model: "claude-3-haiku-20240307", want: "deepseek-chat"},
		{name: "gemini latest alias", model: "gemini-flash-latest", want: "deepseek-chat"},
		{name: "gemini historical pro", model: "gemini-1.5-pro", want: "deepseek-expert-reasoner"},
		{name: "gemini vision legacy", model: "gemini-pro-vision", want: "deepseek-vision"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := ResolveModel(nil, tc.model)
			if !ok || got != tc.want {
				t.Fatalf("expected alias %s -> %s, got ok=%v model=%q", tc.model, tc.want, ok, got)
			}
		})
	}
}

func TestResolveModelUnknown(t *testing.T) {
	_, ok := ResolveModel(nil, "totally-custom-model")
	if ok {
		t.Fatal("expected unknown model to fail resolve")
	}
}

func TestResolveModelUnknownKnownFamilyName(t *testing.T) {
	_, ok := ResolveModel(nil, "gpt-5.5-pro-search")
	if ok {
		t.Fatal("expected unknown known-family model to fail resolve without alias")
	}
}

func TestResolveModelRejectsOldDeepSeekV4IDs(t *testing.T) {
	oldModels := []string{
		"deepseek-v4-flash",
		"deepseek-v4-pro",
		"deepseek-v4-flash-search",
		"deepseek-v4-pro-search",
		"deepseek-v4-vision",
	}
	for _, model := range oldModels {
		if got, ok := ResolveModel(nil, model); ok {
			t.Fatalf("expected old model %q to be rejected, got %q", model, got)
		}
	}
}

func TestResolveModelRejectsRetiredHistoricalModels(t *testing.T) {
	retiredModels := []string{
		"claude-2.1",
		"claude-instant-1.2",
		"gpt-3.5-turbo",
	}
	for _, model := range retiredModels {
		if got, ok := ResolveModel(nil, model); ok {
			t.Fatalf("expected retired model %q to be rejected, got %q", model, got)
		}
	}
}

func TestResolveModelDirectDeepSeekExpert(t *testing.T) {
	got, ok := ResolveModel(nil, "deepseek-expert-reasoner")
	if !ok || got != "deepseek-expert-reasoner" {
		t.Fatalf("expected deepseek-expert-reasoner, got ok=%v model=%q", ok, got)
	}
}

func TestResolveModelCustomAliasToExpert(t *testing.T) {
	got, ok := ResolveModel(mockModelAliasReader{
		"my-expert-model": "deepseek-expert-reasoner-search",
	}, "my-expert-model")
	if !ok || got != "deepseek-expert-reasoner-search" {
		t.Fatalf("expected alias -> deepseek-expert-reasoner-search, got ok=%v model=%q", ok, got)
	}
}

func TestResolveModelCustomAliasToVision(t *testing.T) {
	got, ok := ResolveModel(mockModelAliasReader{
		"my-vision-model": "deepseek-vision",
	}, "my-vision-model")
	if !ok || got != "deepseek-vision" {
		t.Fatalf("expected alias -> deepseek-vision, got ok=%v model=%q", ok, got)
	}
}

func TestClaudeModelsResponsePaginationFields(t *testing.T) {
	resp := ClaudeModelsResponse()
	if _, ok := resp["first_id"]; !ok {
		t.Fatalf("expected first_id in response: %#v", resp)
	}
	if _, ok := resp["last_id"]; !ok {
		t.Fatalf("expected last_id in response: %#v", resp)
	}
	if _, ok := resp["has_more"]; !ok {
		t.Fatalf("expected has_more in response: %#v", resp)
	}
}
