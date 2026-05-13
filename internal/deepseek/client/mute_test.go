package client

import "testing"

func TestExtractMuteInfoFromBizCodeResponse(t *testing.T) {
	info := extractMuteInfo(map[string]any{
		"code": 0,
		"msg":  "",
		"data": map[string]any{
			"biz_code": 5,
			"biz_msg":  "user is muted",
			"biz_data": map[string]any{
				"is_muted":   1,
				"mute_until": 1778706012.595,
			},
		},
	})

	if !info.Muted {
		t.Fatal("expected muted response")
	}
	if info.Until != 1778706012.595 {
		t.Fatalf("unexpected mute until %v", info.Until)
	}
}
