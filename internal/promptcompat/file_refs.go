package promptcompat

import "strings"

const openAIFileTagPrefix = "<||file:"

func CollectOpenAIRefFileIDs(req map[string]any) []string {
	if len(req) == 0 {
		return nil
	}
	out := make([]string, 0, 4)
	seen := map[string]struct{}{}
	for _, key := range []string{
		"ref_file_ids",
		"file_ids",
		"attachments",
		"messages",
		"input",
	} {
		raw := req[key]
		if raw == nil {
			continue
		}
		// Skip top-level strings for 'messages' and 'input' as they are likely plain text content,
		// not file IDs. String file IDs are expected in 'ref_file_ids' or 'file_ids';
		// message/input strings may still contain explicit DS2API file tags.
		if key == "messages" || key == "input" {
			if s, ok := raw.(string); ok {
				appendOpenAIFileTagRefs(&out, seen, s, nil)
				continue
			}
		}
		appendOpenAIRefFileIDs(&out, seen, raw)
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

func appendOpenAIRefFileIDs(out *[]string, seen map[string]struct{}, raw any) {
	switch x := raw.(type) {
	case string:
		addOpenAIRefFileID(out, seen, x)
	case []string:
		for _, item := range x {
			addOpenAIRefFileID(out, seen, item)
		}
	case []any:
		for _, item := range x {
			appendOpenAIRefFileIDs(out, seen, item)
		}
	case map[string]any:
		if fileID := strings.TrimSpace(asString(x["file_id"])); fileID != "" {
			addOpenAIRefFileID(out, seen, fileID)
		}
		if strings.Contains(strings.ToLower(strings.TrimSpace(asString(x["type"]))), "file") {
			if fileID := strings.TrimSpace(asString(x["id"])); fileID != "" {
				addOpenAIRefFileID(out, seen, fileID)
			}
		}
		if fileMap, ok := x["file"].(map[string]any); ok {
			if fileID := strings.TrimSpace(asString(fileMap["file_id"])); fileID != "" {
				addOpenAIRefFileID(out, seen, fileID)
			}
			if fileID := strings.TrimSpace(asString(fileMap["id"])); fileID != "" {
				addOpenAIRefFileID(out, seen, fileID)
			}
		}
		// Recurse into potential containers. Note: we do NOT recurse into 'content' or 'input'
		// if they are plain strings (handled by the top-level switch), but they are usually
		// nested inside the map branch anyway.
		// To be safe, we only recurse into these known container keys.
		for _, key := range []string{"ref_file_ids", "file_ids", "attachments", "messages", "input", "content", "files", "items", "data", "source"} {
			if nested, ok := x[key]; ok {
				// If it's a message content that is a string, we must NOT treat it as an ID.
				if key == "content" || key == "input" {
					if s, ok := nested.(string); ok {
						appendOpenAIFileTagRefs(out, seen, s, nil)
						continue
					}
				}
				appendOpenAIRefFileIDs(out, seen, nested)
			}
		}
	}
}

func addOpenAIRefFileID(out *[]string, seen map[string]struct{}, fileID string) {
	fileID = strings.TrimSpace(fileID)
	if fileID == "" {
		return
	}
	if _, ok := seen[fileID]; ok {
		return
	}
	seen[fileID] = struct{}{}
	*out = append(*out, fileID)
}

func CollectOpenAIFileTagAccounts(req map[string]any) []string {
	accounts := make([]string, 0, 1)
	seen := map[string]struct{}{}
	collectOpenAIFileTagAccounts(req, &accounts, seen)
	if len(accounts) == 0 {
		return nil
	}
	return accounts
}

func collectOpenAIFileTagAccounts(raw any, out *[]string, seen map[string]struct{}) {
	switch x := raw.(type) {
	case string:
		appendOpenAIFileTagRefs(nil, nil, x, func(account string) {
			account = strings.TrimSpace(account)
			if account == "" {
				return
			}
			key := strings.ToLower(account)
			if _, ok := seen[key]; ok {
				return
			}
			seen[key] = struct{}{}
			*out = append(*out, account)
		})
	case []any:
		for _, item := range x {
			collectOpenAIFileTagAccounts(item, out, seen)
		}
	case map[string]any:
		for _, key := range []string{"attachments", "messages", "input", "content", "files", "items", "data", "source"} {
			if nested, ok := x[key]; ok {
				collectOpenAIFileTagAccounts(nested, out, seen)
			}
		}
	}
}

func StripOpenAIFileTags(raw any) any {
	switch x := raw.(type) {
	case string:
		return stripOpenAIFileTagsFromString(x)
	case []any:
		for i, item := range x {
			x[i] = StripOpenAIFileTags(item)
		}
		return x
	case map[string]any:
		for _, key := range []string{"attachments", "messages", "input", "content", "files", "items", "data", "source"} {
			if nested, ok := x[key]; ok {
				x[key] = StripOpenAIFileTags(nested)
			}
		}
		return x
	default:
		return raw
	}
}

func stripOpenAIFileTagsFromString(s string) string {
	if !strings.Contains(s, openAIFileTagPrefix) {
		return s
	}
	var b strings.Builder
	last := 0
	for _, tag := range scanOpenAIFileTags(s) {
		b.WriteString(s[last:tag.start])
		last = tag.end
	}
	if last == 0 {
		return s
	}
	b.WriteString(s[last:])
	return strings.TrimSpace(b.String())
}

type openAIFileTagRef struct {
	fileID  string
	account string
	start   int
	end     int
}

func appendOpenAIFileTagRefs(out *[]string, seen map[string]struct{}, s string, onAccount func(string)) {
	for _, tag := range scanOpenAIFileTags(s) {
		if out != nil {
			addOpenAIRefFileID(out, seen, tag.fileID)
		}
		if onAccount != nil {
			onAccount(tag.account)
		}
	}
}

func scanOpenAIFileTags(s string) []openAIFileTagRef {
	if !strings.Contains(s, openAIFileTagPrefix) {
		return nil
	}
	refs := make([]openAIFileTagRef, 0, 1)
	searchFrom := 0
	for {
		relStart := strings.Index(s[searchFrom:], openAIFileTagPrefix)
		if relStart < 0 {
			break
		}
		start := searchFrom + relStart
		contentStart := start + len(openAIFileTagPrefix)
		relEnd := strings.Index(s[contentStart:], "||>")
		if relEnd < 0 {
			break
		}
		end := contentStart + relEnd + len("||>")
		body := s[contentStart : contentStart+relEnd]
		if ref, ok := parseOpenAIFileTagBody(body); ok {
			ref.start = start
			ref.end = end
			refs = append(refs, ref)
		}
		searchFrom = end
	}
	return refs
}

func parseOpenAIFileTagBody(body string) (openAIFileTagRef, bool) {
	lastColon := strings.LastIndex(body, ":")
	if lastColon <= 0 || lastColon >= len(body)-1 {
		return openAIFileTagRef{}, false
	}
	account := strings.TrimSpace(body[lastColon+1:])
	rest := body[:lastColon]
	secondLastColon := strings.LastIndex(rest, ":")
	if secondLastColon <= 0 || secondLastColon >= len(rest)-1 {
		return openAIFileTagRef{}, false
	}
	fileID := strings.TrimSpace(rest[secondLastColon+1:])
	if fileID == "" || account == "" {
		return openAIFileTagRef{}, false
	}
	return openAIFileTagRef{fileID: fileID, account: account}, true
}
