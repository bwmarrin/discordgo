package discordgo

import (
	"bytes"
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestVerifyInteraction(t *testing.T) {
	pubkey, privkey, err := ed25519.GenerateKey(nil)
	if err != nil {
		t.Errorf("error generating signing keypair: %s", err)
	}
	timestamp := "1608597133"

	t.Run("success", func(t *testing.T) {
		body := "body"
		request := httptest.NewRequest("POST", "http://localhost/interaction", strings.NewReader(body))
		request.Header.Set("X-Signature-Timestamp", timestamp)

		var msg bytes.Buffer
		msg.WriteString(timestamp)
		msg.WriteString(body)
		signature := ed25519.Sign(privkey, msg.Bytes())
		request.Header.Set("X-Signature-Ed25519", hex.EncodeToString(signature[:ed25519.SignatureSize]))

		if !VerifyInteraction(request, pubkey) {
			t.Error("expected true, got false")
		}
	})

	t.Run("failure/modified body", func(t *testing.T) {
		body := "body"
		request := httptest.NewRequest("POST", "http://localhost/interaction", strings.NewReader("WRONG"))
		request.Header.Set("X-Signature-Timestamp", timestamp)

		var msg bytes.Buffer
		msg.WriteString(timestamp)
		msg.WriteString(body)
		signature := ed25519.Sign(privkey, msg.Bytes())
		request.Header.Set("X-Signature-Ed25519", hex.EncodeToString(signature[:ed25519.SignatureSize]))

		if VerifyInteraction(request, pubkey) {
			t.Error("expected false, got true")
		}
	})

	t.Run("failure/modified timestamp", func(t *testing.T) {
		body := "body"
		request := httptest.NewRequest("POST", "http://localhost/interaction", strings.NewReader("WRONG"))
		request.Header.Set("X-Signature-Timestamp", strconv.FormatInt(time.Now().Add(time.Minute).Unix(), 10))

		var msg bytes.Buffer
		msg.WriteString(timestamp)
		msg.WriteString(body)
		signature := ed25519.Sign(privkey, msg.Bytes())
		request.Header.Set("X-Signature-Ed25519", hex.EncodeToString(signature[:ed25519.SignatureSize]))

		if VerifyInteraction(request, pubkey) {
			t.Error("expected false, got true")
		}
	})
}

func TestModalSubmitRoleSelect(t *testing.T) {
	raw := []byte(`{
		"id": "1111111111111111111",
		"application_id": "2222222222222222222",
		"type": 5,
		"token": "test-token",
		"version": 1,
		"guild_id": "3333333333333333333",
		"channel_id": "4444444444444444444",
		"data": {
			"custom_id": "help_admin_panel",
			"components": [
				{
					"type": 1,
					"components": [
						{
							"type": 6,
							"custom_id": "role_select",
							"values": ["8888888888888888888"]
						}
					]
				}
			],
			"resolved": {
				"roles": {
					"8888888888888888888": {
						"id": "8888888888888888888",
						"name": "Admin",
						"color": 16711680,
						"hoist": true,
						"position": 1,
						"permissions": "8",
						"managed": false,
						"mentionable": true
					}
				}
			}
		}
	}`)

	var i Interaction
	if err := json.Unmarshal(raw, &i); err != nil {
		t.Fatalf("unexpected unmarshal error: %v", err)
	}

	data := i.ModalSubmitData()

	if data.CustomID != "help_admin_panel" {
		t.Errorf("expected custom_id %q, got %q", "help_admin_panel", data.CustomID)
	}

	if len(data.Components) != 1 {
		t.Fatalf("expected 1 top-level component, got %d", len(data.Components))
	}
	row, ok := data.Components[0].(*ActionsRow)
	if !ok {
		t.Fatalf("expected *ActionsRow at top level, got %T", data.Components[0])
	}
	if len(row.Components) != 1 {
		t.Fatalf("expected 1 inner component, got %d", len(row.Components))
	}
	sel, ok := row.Components[0].(*SelectMenu)
	if !ok {
		t.Fatalf("expected *SelectMenu inner component, got %T", row.Components[0])
	}
	if sel.CustomID != "role_select" {
		t.Errorf("expected custom_id %q, got %q", "role_select", sel.CustomID)
	}
	if sel.MenuType != RoleSelectMenu {
		t.Errorf("expected MenuType %d (RoleSelectMenu), got %d", RoleSelectMenu, sel.MenuType)
	}
	if len(sel.Values) != 1 || sel.Values[0] != "8888888888888888888" {
		t.Errorf("expected values [8888888888888888888], got %v", sel.Values)
	}

	if data.Resolved.Roles == nil {
		t.Fatal("expected resolved roles map to be non-nil")
	}
	role, exists := data.Resolved.Roles["8888888888888888888"]
	if !exists {
		t.Fatal("expected resolved role 8888888888888888888 to exist")
	}
	if role.Name != "Admin" {
		t.Errorf("expected role name %q, got %q", "Admin", role.Name)
	}
	if role.Color != 16711680 {
		t.Errorf("expected role color 16711680, got %d", role.Color)
	}

	found := data.GetSelectMenu("role_select")
	if found == nil {
		t.Fatal("GetSelectMenu returned nil for existing custom_id")
	}
	if found != sel {
		t.Error("GetSelectMenu returned a different *SelectMenu than the one in the tree")
	}
}

func TestModalSubmitChannelSelect(t *testing.T) {
	raw := []byte(`{
		"id": "1111111111111111112",
		"application_id": "2222222222222222222",
		"type": 5,
		"token": "test-token",
		"version": 1,
		"guild_id": "3333333333333333333",
		"channel_id": "4444444444444444444",
		"data": {
			"custom_id": "channel_config_modal",
			"components": [
				{
					"type": 1,
					"components": [
						{
							"type": 8,
							"custom_id": "channel_select",
							"values": ["9999999999999999999"]
						}
					]
				}
			],
			"resolved": {
				"channels": {
					"9999999999999999999": {
						"id": "9999999999999999999",
						"type": 0,
						"name": "general",
						"permissions": "1071698529857"
					}
				}
			}
		}
	}`)

	var i Interaction
	if err := json.Unmarshal(raw, &i); err != nil {
		t.Fatalf("unexpected unmarshal error: %v", err)
	}

	data := i.ModalSubmitData()

	sel := data.GetSelectMenu("channel_select")
	if sel == nil {
		t.Fatal("GetSelectMenu returned nil for channel_select")
	}
	if sel.MenuType != ChannelSelectMenu {
		t.Errorf("expected MenuType %d (ChannelSelectMenu), got %d", ChannelSelectMenu, sel.MenuType)
	}
	if len(sel.Values) != 1 || sel.Values[0] != "9999999999999999999" {
		t.Errorf("expected values [9999999999999999999], got %v", sel.Values)
	}

	ch, exists := data.Resolved.Channels["9999999999999999999"]
	if !exists {
		t.Fatal("expected resolved channel 9999999999999999999 to exist")
	}
	if ch.Name != "general" {
		t.Errorf("expected channel name %q, got %q", "general", ch.Name)
	}
}

func TestModalSubmitUserSelect(t *testing.T) {
	raw := []byte(`{
		"id": "1111111111111111113",
		"application_id": "2222222222222222222",
		"type": 5,
		"token": "test-token",
		"version": 1,
		"guild_id": "3333333333333333333",
		"channel_id": "4444444444444444444",
		"data": {
			"custom_id": "user_assign_modal",
			"components": [
				{
					"type": 1,
					"components": [
						{
							"type": 5,
							"custom_id": "user_select",
							"values": ["123456789012345678"]
						}
					]
				}
			],
			"resolved": {
				"users": {
					"123456789012345678": {
						"id": "123456789012345678",
						"username": "octocat",
						"discriminator": "0001"
					}
				},
				"members": {
					"123456789012345678": {
						"nick": "GitHub Octocat",
						"roles": ["8888888888888888888"]
					}
				}
			}
		}
	}`)

	var i Interaction
	if err := json.Unmarshal(raw, &i); err != nil {
		t.Fatalf("unexpected unmarshal error: %v", err)
	}

	data := i.ModalSubmitData()

	sel := data.GetSelectMenu("user_select")
	if sel == nil {
		t.Fatal("GetSelectMenu returned nil for user_select")
	}
	if sel.MenuType != UserSelectMenu {
		t.Errorf("expected MenuType %d (UserSelectMenu), got %d", UserSelectMenu, sel.MenuType)
	}
	if len(sel.Values) != 1 || sel.Values[0] != "123456789012345678" {
		t.Errorf("expected values [123456789012345678], got %v", sel.Values)
	}

	u, exists := data.Resolved.Users["123456789012345678"]
	if !exists {
		t.Fatal("expected resolved user 123456789012345678 to exist")
	}
	if u.Username != "octocat" {
		t.Errorf("expected username %q, got %q", "octocat", u.Username)
	}

	m, exists := data.Resolved.Members["123456789012345678"]
	if !exists {
		t.Fatal("expected resolved member 123456789012345678 to exist")
	}
	if m.Nick != "GitHub Octocat" {
		t.Errorf("expected member nick %q, got %q", "GitHub Octocat", m.Nick)
	}
}

func TestModalSubmitMentionableSelect(t *testing.T) {
	raw := []byte(`{
		"id": "1111111111111111114",
		"application_id": "2222222222222222222",
		"type": 5,
		"token": "test-token",
		"version": 1,
		"data": {
			"custom_id": "mentionable_modal",
			"components": [
				{
					"type": 1,
					"components": [
						{
							"type": 7,
							"custom_id": "mentionable_select",
							"values": ["8888888888888888888", "123456789012345678"]
						}
					]
				}
			],
			"resolved": {
				"roles": {
					"8888888888888888888": {
						"id": "8888888888888888888",
						"name": "Moderators"
					}
				},
				"users": {
					"123456789012345678": {
						"id": "123456789012345678",
						"username": "mod_user"
					}
				}
			}
		}
	}`)

	var i Interaction
	if err := json.Unmarshal(raw, &i); err != nil {
		t.Fatalf("unexpected unmarshal error: %v", err)
	}

	data := i.ModalSubmitData()

	sel := data.GetSelectMenu("mentionable_select")
	if sel == nil {
		t.Fatal("GetSelectMenu returned nil for mentionable_select")
	}
	if sel.MenuType != MentionableSelectMenu {
		t.Errorf("expected MenuType %d (MentionableSelectMenu), got %d", MentionableSelectMenu, sel.MenuType)
	}
	if len(sel.Values) != 2 {
		t.Fatalf("expected 2 values, got %d", len(sel.Values))
	}

	if _, ok := data.Resolved.Roles["8888888888888888888"]; !ok {
		t.Error("expected resolved role 8888888888888888888 to exist")
	}
	if _, ok := data.Resolved.Users["123456789012345678"]; !ok {
		t.Error("expected resolved user 123456789012345678 to exist")
	}
}

func TestModalSubmitStringSelect(t *testing.T) {
	raw := []byte(`{
		"id": "1111111111111111115",
		"application_id": "2222222222222222222",
		"type": 5,
		"token": "test-token",
		"version": 1,
		"data": {
			"custom_id": "string_select_modal",
			"components": [
				{
					"type": 1,
					"components": [
						{
							"type": 3,
							"custom_id": "rating_select",
							"values": ["awesome"]
						}
					]
				}
			]
		}
	}`)

	var i Interaction
	if err := json.Unmarshal(raw, &i); err != nil {
		t.Fatalf("unexpected unmarshal error: %v", err)
	}

	data := i.ModalSubmitData()

	sel := data.GetSelectMenu("rating_select")
	if sel == nil {
		t.Fatal("GetSelectMenu returned nil for rating_select")
	}
	if sel.MenuType != StringSelectMenu {
		t.Errorf("expected MenuType %d (StringSelectMenu), got %d", StringSelectMenu, sel.MenuType)
	}
	if len(sel.Values) != 1 || sel.Values[0] != "awesome" {
		t.Errorf("expected values [awesome], got %v", sel.Values)
	}
}

func TestModalSubmitTextInput(t *testing.T) {
	raw := []byte(`{
		"id": "1111111111111111116",
		"application_id": "2222222222222222222",
		"type": 5,
		"token": "test-token",
		"version": 1,
		"guild_id": "3333333333333333333",
		"channel_id": "4444444444444444444",
		"data": {
			"custom_id": "text_modal",
			"components": [
				{
					"type": 1,
					"components": [
						{
							"type": 4,
							"custom_id": "my_input",
							"value": "hello world",
							"style": 1,
							"label": "Name"
						}
					]
				}
			]
		}
	}`)

	var i Interaction
	if err := json.Unmarshal(raw, &i); err != nil {
		t.Fatalf("unexpected unmarshal error: %v", err)
	}

	data := i.ModalSubmitData()

	if data.CustomID != "text_modal" {
		t.Errorf("expected custom_id %q, got %q", "text_modal", data.CustomID)
	}

	ti := data.GetTextInput("my_input")
	if ti == nil {
		t.Fatal("GetTextInput returned nil for my_input")
	}
	if ti.Value != "hello world" {
		t.Errorf("expected value %q, got %q", "hello world", ti.Value)
	}

	if sel := data.GetSelectMenu("my_input"); sel != nil {
		t.Error("GetSelectMenu should return nil for a TextInput custom_id")
	}
}

func TestModalSubmitMultipleSelectMenus(t *testing.T) {
	raw := []byte(`{
		"id": "1111111111111111117",
		"application_id": "2222222222222222222",
		"type": 5,
		"token": "test-token",
		"version": 1,
		"data": {
			"custom_id": "multi_select_modal",
			"components": [
				{
					"type": 1,
					"components": [
						{
							"type": 6,
							"custom_id": "role_picker",
							"values": ["1234000000000000000"]
						}
					]
				},
				{
					"type": 1,
					"components": [
						{
							"type": 8,
							"custom_id": "channel_picker",
							"values": ["5678000000000000000"]
						}
					]
				}
			],
			"resolved": {
				"roles": {
					"1234000000000000000": {
						"id": "1234000000000000000",
						"name": "Moderator",
						"color": 0,
						"hoist": false,
						"position": 2,
						"permissions": "0",
						"managed": false,
						"mentionable": false
					}
				},
				"channels": {
					"5678000000000000000": {
						"id": "5678000000000000000",
						"type": 0,
						"name": "announcements",
						"permissions": "0"
					}
				}
			}
		}
	}`)

	var i Interaction
	if err := json.Unmarshal(raw, &i); err != nil {
		t.Fatalf("unexpected unmarshal error: %v", err)
	}

	data := i.ModalSubmitData()

	roleSel := data.GetSelectMenu("role_picker")
	if roleSel == nil {
		t.Fatal("GetSelectMenu returned nil for role_picker")
	}
	if roleSel.MenuType != RoleSelectMenu {
		t.Errorf("expected RoleSelectMenu, got %d", roleSel.MenuType)
	}
	if len(roleSel.Values) != 1 || roleSel.Values[0] != "1234000000000000000" {
		t.Errorf("unexpected role values: %v", roleSel.Values)
	}

	chanSel := data.GetSelectMenu("channel_picker")
	if chanSel == nil {
		t.Fatal("GetSelectMenu returned nil for channel_picker")
	}
	if chanSel.MenuType != ChannelSelectMenu {
		t.Errorf("expected ChannelSelectMenu, got %d", chanSel.MenuType)
	}
	if len(chanSel.Values) != 1 || chanSel.Values[0] != "5678000000000000000" {
		t.Errorf("unexpected channel values: %v", chanSel.Values)
	}

	if _, ok := data.Resolved.Roles["1234000000000000000"]; !ok {
		t.Error("expected resolved role 1234000000000000000 to exist")
	}
	if _, ok := data.Resolved.Channels["5678000000000000000"]; !ok {
		t.Error("expected resolved channel 5678000000000000000 to exist")
	}
}

func TestModalSubmitSelectInLabel(t *testing.T) {
	raw := []byte(`{
		"id": "1111111111111111118",
		"application_id": "2222222222222222222",
		"type": 5,
		"token": "test-token",
		"version": 1,
		"data": {
			"custom_id": "label_modal",
			"components": [
				{
					"type": 18,
					"label": "Pick a role",
					"component": {
						"type": 6,
						"custom_id": "labeled_role_select",
						"values": ["7777777777777777777"]
					}
				}
			],
			"resolved": {
				"roles": {
					"7777777777777777777": {
						"id": "7777777777777777777",
						"name": "Member",
						"color": 255,
						"hoist": false,
						"position": 0,
						"permissions": "0",
						"managed": false,
						"mentionable": true
					}
				}
			}
		}
	}`)

	var i Interaction
	if err := json.Unmarshal(raw, &i); err != nil {
		t.Fatalf("unexpected unmarshal error: %v", err)
	}

	data := i.ModalSubmitData()

	if len(data.Components) != 1 {
		t.Fatalf("expected 1 top-level component, got %d", len(data.Components))
	}
	lbl, ok := data.Components[0].(*Label)
	if !ok {
		t.Fatalf("expected *Label at top level, got %T", data.Components[0])
	}
	if lbl.Label != "Pick a role" {
		t.Errorf("expected label %q, got %q", "Pick a role", lbl.Label)
	}

	sel := data.GetSelectMenu("labeled_role_select")
	if sel == nil {
		t.Fatal("GetSelectMenu returned nil for select inside Label")
	}
	if sel.MenuType != RoleSelectMenu {
		t.Errorf("expected RoleSelectMenu, got %d", sel.MenuType)
	}
	if len(sel.Values) != 1 || sel.Values[0] != "7777777777777777777" {
		t.Errorf("unexpected values: %v", sel.Values)
	}

	if _, ok := data.Resolved.Roles["7777777777777777777"]; !ok {
		t.Error("expected resolved role 7777777777777777777 to exist")
	}
}

func TestModalSubmitTextInputInLabel(t *testing.T) {
	raw := []byte(`{
		"id": "1111111111111111119",
		"application_id": "2222222222222222222",
		"type": 5,
		"token": "test-token",
		"version": 1,
		"data": {
			"custom_id": "label_text_modal",
			"components": [
				{
					"type": 18,
					"label": "Suggestions",
					"component": {
						"type": 4,
						"custom_id": "labeled_text_input",
						"value": "Great experience!"
					}
				}
			]
		}
	}`)

	var i Interaction
	if err := json.Unmarshal(raw, &i); err != nil {
		t.Fatalf("unexpected unmarshal error: %v", err)
	}

	data := i.ModalSubmitData()

	if len(data.Components) != 1 {
		t.Fatalf("expected 1 top-level component, got %d", len(data.Components))
	}
	lbl, ok := data.Components[0].(*Label)
	if !ok {
		t.Fatalf("expected *Label at top level, got %T", data.Components[0])
	}
	if lbl.Label != "Suggestions" {
		t.Errorf("expected label %q, got %q", "Suggestions", lbl.Label)
	}

	ti := data.GetTextInput("labeled_text_input")
	if ti == nil {
		t.Fatal("GetTextInput returned nil for text input inside Label")
	}
	if ti.Value != "Great experience!" {
		t.Errorf("expected value %q, got %q", "Great experience!", ti.Value)
	}
}

func TestModalSubmitGetMissing(t *testing.T) {
	data := ModalSubmitInteractionData{
		CustomID: "empty_modal",
		Components: []MessageComponent{
			&ActionsRow{
				Components: []MessageComponent{
					&SelectMenu{MenuType: RoleSelectMenu, CustomID: "role_select"},
					&TextInput{CustomID: "text_input", Value: "val"},
				},
			},
			ActionsRow{
				Components: []MessageComponent{
					&SelectMenu{MenuType: ChannelSelectMenu, CustomID: "channel_select_val"},
				},
			},
			&Label{
				Label: "Label 1",
				Component: &SelectMenu{
					MenuType: UserSelectMenu,
					CustomID: "user_select_in_label",
				},
			},
			Label{
				Label: "Label 2",
				Component: &TextInput{
					CustomID: "text_in_label_val",
					Value:    "hello",
				},
			},
			&SelectMenu{MenuType: StringSelectMenu, CustomID: "top_level_select"},
			&TextInput{CustomID: "top_level_input", Value: "top"},
		},
	}

	if sel := data.GetSelectMenu("nonexistent"); sel != nil {
		t.Error("expected nil for nonexistent custom_id in GetSelectMenu, got a result")
	}
	if ti := data.GetTextInput("nonexistent"); ti != nil {
		t.Error("expected nil for nonexistent custom_id in GetTextInput, got a result")
	}

	if sel := data.GetSelectMenu("role_select"); sel == nil {
		t.Error("expected non-nil for existing custom_id role_select")
	}
	if sel := data.GetSelectMenu("channel_select_val"); sel == nil {
		t.Error("expected non-nil for value ActionsRow channel_select_val")
	}
	if sel := data.GetSelectMenu("user_select_in_label"); sel == nil {
		t.Error("expected non-nil for user_select_in_label")
	}
	if sel := data.GetSelectMenu("top_level_select"); sel == nil {
		t.Error("expected non-nil for top_level_select")
	}

	if ti := data.GetTextInput("text_input"); ti == nil {
		t.Error("expected non-nil for existing custom_id text_input")
	}
	if ti := data.GetTextInput("text_in_label_val"); ti == nil {
		t.Error("expected non-nil for value Label text_in_label_val")
	}
	if ti := data.GetTextInput("top_level_input"); ti == nil {
		t.Error("expected non-nil for top_level_input")
	}
}
