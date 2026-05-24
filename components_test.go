package discordgo

import (
	"encoding/json"
	"testing"
)

func TestMessageCreateUnknownComponentType(t *testing.T) {
	var m MessageCreate
	err := json.Unmarshal([]byte(`{
		"id":"message",
		"channel_id":"channel",
		"content":"content",
		"components":[{"type":20,"id":1,"custom":"value"}]
	}`), &m)
	if err != nil {
		t.Fatalf("json.Unmarshal returned error: %v", err)
	}
	if m.Message == nil {
		t.Fatal("Message is nil")
	}
	if len(m.Components) != 1 {
		t.Fatalf("len(Components) = %d, want 1", len(m.Components))
	}
	if m.Components[0].Type() != ComponentType(20) {
		t.Fatalf("component type = %d, want 20", m.Components[0].Type())
	}

	raw, err := m.Components[0].MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON returned error: %v", err)
	}

	var component struct {
		Type   ComponentType `json:"type"`
		ID     int           `json:"id"`
		Custom string        `json:"custom"`
	}
	if err := json.Unmarshal(raw, &component); err != nil {
		t.Fatalf("json.Unmarshal component returned error: %v", err)
	}
	if component.Type != ComponentType(20) || component.ID != 1 || component.Custom != "value" {
		t.Fatalf("component = %#v, want type 20 with original fields", component)
	}
}
