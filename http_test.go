package discordgo

import (
	"bytes"
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

func signRequest(key ed25519.PrivateKey, r *http.Request) error {
	timestamp := time.Now().Unix()
	r.Header.Set("X-Signature-Timestamp", fmt.Sprint(timestamp))
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("failed to read request body. err: %w", err)
	}
	var msg bytes.Buffer
	msg.WriteString(fmt.Sprint(timestamp))
	msg.WriteString(string(body))

	signature := ed25519.Sign(key, msg.Bytes())
	r.Header.Set("X-Signature-Ed25519", hex.EncodeToString(signature[:ed25519.SignatureSize]))

	r.Body = io.NopCloser(bytes.NewBuffer(body))
	return nil
}

func TestServeHTTP(t *testing.T) {
	pubkey, privkey, err := ed25519.GenerateKey(nil)
	if err != nil {
		t.Errorf("error generating signing keypair: %s", err)
	}
	mux := NewInteractionsHTTPServer(dg, hex.EncodeToString(pubkey))

	testHandlerCalled := int32(0)
	testHandler := func(s *Session, ic *InteractionCreate) {
		atomic.AddInt32(&testHandlerCalled, 1)
	}
	mux.HandleFunc("basic", testHandler)
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(InteractionCreate{
		&Interaction{ID: "fake",
			Type: InteractionApplicationCommand,
			Data: ApplicationCommandInteractionData{Name: "basic"},
		}})
	request, _ := http.NewRequest(http.MethodPost, "/", b)
	response := httptest.NewRecorder()
	err = signRequest(privkey, request)
	if err != nil {
		t.Fatalf("failed to sign request. err: %v", err)
	}

	mux.ServeHTTP(response, request)

	if atomic.LoadInt32(&testHandlerCalled) != 1 {
		t.Fatalf("testHandler was not called once.")
	}
}
