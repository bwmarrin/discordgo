package discordgo

import (
	"bytes"
	"crypto/ed25519"
	"encoding/hex"
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

// TestVerifyUnmarshalInteraction tests the combined Verify and Unmarshal for interactions.
func TestVerifyUnmarshalInteraction(t *testing.T) {
	pubkey, privkey, err := ed25519.GenerateKey(nil)
	if err != nil {
		t.Fatalf("error generating signing keypair: %s", err)
	}
	timestamp := "1608597133"
	// prepare a minimal valid interaction JSON (ApplicationCommand)
	payload := `{
       "id": "123",
       "application_id": "appid",
       "type": 2,
       "data": {"id":"cmdid","name":"test","type":1}
    }`

	t.Run("success", func(t *testing.T) {
		body := payload
		req := httptest.NewRequest("POST", "http://example.com/", strings.NewReader(body))
		req.Header.Set("X-Signature-Timestamp", timestamp)

		var msg bytes.Buffer
		msg.WriteString(timestamp)
		msg.WriteString(body)

		sig := ed25519.Sign(privkey, msg.Bytes())
		req.Header.Set("X-Signature-Ed25519", hex.EncodeToString(sig))
		var inter Interaction
		if err := VerifyUnmarshalInteraction(req, pubkey, &inter); err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if inter.ID != "123" || inter.Type != InteractionApplicationCommand {
			t.Errorf("unexpected interaction data: %+v", inter)
		}
	})

	t.Run("missing signature header", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/", strings.NewReader(payload))
		req.Header.Set("X-Signature-Timestamp", timestamp)
		var inter Interaction
		err := VerifyUnmarshalInteraction(req, pubkey, &inter)
		if err == nil || !strings.Contains(err.Error(), "Missing X-Signature-Ed25519 header") {
			t.Errorf("expected missing signature error, got %v", err)
		}
	})

	t.Run("invalid signature hex", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/", strings.NewReader(payload))
		req.Header.Set("X-Signature-Timestamp", timestamp)
		req.Header.Set("X-Signature-Ed25519", "zzzz")
		var inter Interaction
		err := VerifyUnmarshalInteraction(req, pubkey, &inter)
		if err == nil || !strings.Contains(err.Error(), "Failed to decode Ed25519 signature") {
			t.Errorf("expected decode error, got %v", err)
		}
	})

	t.Run("invalid signature size", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/", strings.NewReader(payload))
		req.Header.Set("X-Signature-Timestamp", timestamp)
		// too short sig
		short := make([]byte, ed25519.SignatureSize/2)
		req.Header.Set("X-Signature-Ed25519", hex.EncodeToString(short))
		var inter Interaction
		err := VerifyUnmarshalInteraction(req, pubkey, &inter)
		if err == nil || !strings.Contains(err.Error(), "Incorrect Ed25519 signature size") {
			t.Errorf("expected signature size error, got %v", err)
		}
	})

	t.Run("missing timestamp header", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/", strings.NewReader(payload))
		sig := ed25519.Sign(privkey, []byte(timestamp+payload))
		req.Header.Set("X-Signature-Ed25519", hex.EncodeToString(sig))
		var inter Interaction
		err := VerifyUnmarshalInteraction(req, pubkey, &inter)
		if err == nil || !strings.Contains(err.Error(), "Missing X-Signature-Timestamp header") {
			t.Errorf("expected missing timestamp error, got %v", err)
		}
	})

	t.Run("verification failure", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/", strings.NewReader(payload))
		req.Header.Set("X-Signature-Timestamp", timestamp)
		// sign wrong message
		sig := ed25519.Sign(privkey, []byte(timestamp+"wrong"))
		req.Header.Set("X-Signature-Ed25519", hex.EncodeToString(sig))
		var inter Interaction
		err := VerifyUnmarshalInteraction(req, pubkey, &inter)
		if err == nil || !strings.Contains(err.Error(), "Ed25519 signature verification failed") {
			t.Errorf("expected verification failure, got %v", err)
		}
	})

	t.Run("invalid json body", func(t *testing.T) {
		body := "not-json"
		req := httptest.NewRequest("POST", "/", strings.NewReader(body))
		req.Header.Set("X-Signature-Timestamp", timestamp)

		var msg bytes.Buffer
		msg.WriteString(timestamp)
		msg.WriteString(body)
		sig := ed25519.Sign(privkey, msg.Bytes())

		req.Header.Set("X-Signature-Ed25519", hex.EncodeToString(sig))
		var inter Interaction
		err := VerifyUnmarshalInteraction(req, pubkey, &inter)
		if err == nil || !strings.Contains(err.Error(), "Failed to unmarshal interaction") {
			t.Errorf("expected unmarshal error, got %v", err)
		}
	})
}
