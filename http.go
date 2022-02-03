package discordgo

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type WebhookServer struct {
	sess      *Session
	publicKey ed25519.PublicKey
}

func NewWebhookServer(sess *Session, pubKeyString string) *WebhookServer {
	key, err := hex.DecodeString(pubKeyString)
	if err != nil {
		log.Fatal("couldn't parse public key string")
	}
	return &WebhookServer{sess: sess, publicKey: key}
}

func (s *WebhookServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//	if ok := VerifyInteraction(r, s.publicKey); !ok {
	//		http.Error(w, "invalid request signature", http.StatusUnauthorized)
	//		return
	//	}
	var interaction Interaction
	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("could not read request body: %v", err)
		http.Error(w, "could not read request body", http.StatusInternalServerError)
		return
	}
	err = interaction.UnmarshalJSON(data)
	if err != nil {
		http.Error(w, "could not parse interaction", http.StatusBadRequest)
		log.Printf("could not parse request interaction: %v", err)
		return
	}
	if interaction.Type == InteractionPing {
		response := InteractionResponse{Type: InteractionResponsePong}

		b, err := json.Marshal(response)
		if err != nil {
			log.Printf("could not marshal response: %v", err)
			return
		}
		w.Write(b)
		return
	}
	for _, eh := range s.sess.handlers[interaction.Type.String()] {
		go eh.eventHandler.Handle(s.sess, interaction)
	}

}
