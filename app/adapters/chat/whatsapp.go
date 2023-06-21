package chat_adapter

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/62teknologi/62sailfish/config"
	"github.com/google/martian/log"
)

// VonageWhatsApp represents the implementation of ChatService for VonageWhatsApp
type VonageWhatsApp struct {
	config config.Config
}

func NewVonageWhatsapp(config config.Config) *VonageWhatsApp {
	return &VonageWhatsApp{
		config: config,
	}
}

type VonageWhatsAppMessage struct {
	From        string `json:"from"`
	To          string `json:"to"`
	MessageType string `json:"message_type"`
	Text        string `json:"text"`
	Channel     string `json:"channel"`
}

func (w *VonageWhatsApp) SendMessage(sender string, recipient string, message string) error {
	if sender == "" {
		sender = w.config.VonageWhatsAppDefaultSender
	}

	url := w.config.VonageWhatsAppUrl
	username := w.config.VonageUsername
	password := w.config.VonagePassword

	// Create the payload struct
	payload := VonageWhatsAppMessage{
		From:        sender,
		To:          recipient,
		MessageType: "text",
		Text:        message,
		Channel:     "whatsapp",
	}

	// Convert the payload to JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return err
	}

	// Set the necessary headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Add Basic Authentication header
	auth := username + ":" + password
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(auth)))

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Errorf("error while close body")
		}
	}(resp.Body)

	// Check the response
	if (resp.StatusCode != http.StatusOK) && (resp.StatusCode != http.StatusAccepted) {
		return fmt.Errorf("request failed with status: %v", resp.Status)
	}

	fmt.Println("VonageWhatsApp message sent successfully!")
	return nil
}
