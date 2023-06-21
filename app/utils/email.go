package utils

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"os"

	"github.com/62teknologi/62sailfish/config"

	"gopkg.in/gomail.v2"
)

type EmailReceiver struct {
	Subject string
	Name    string
	Address string
}

func EmailSender(htmlTemplate string, params any, receiverList []EmailReceiver) {
	loadedConfig, err := config.LoadConfig(".")
	if err != nil {
		fmt.Errorf("error while load config: %w", err)
		return
	}

	d := gomail.NewDialer(loadedConfig.EmailSMTPHost, loadedConfig.EmailSMTPPort, loadedConfig.EmailAUTHUsername, loadedConfig.EmailAUTHPassword)
	s, err := d.Dial()
	if err != nil {
		fmt.Errorf("error while setup email config: %w", err)
		return
	}

	templateFile, err := os.Open("public/" + htmlTemplate)
	if err != nil {
		fmt.Errorf("error while load template: %w", err)
		return
	}
	defer templateFile.Close()

	// Get Template
	templateBytes, err := io.ReadAll(templateFile)
	if err != nil {
		fmt.Errorf("error while convert template to string: %w", err)
		return
	}
	templateString := string(templateBytes)

	t, err := template.New("webpage").Parse(templateString)
	if err != nil {
		fmt.Errorf("error while parse template string: %w", err)
		return
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, params)
	if err != nil {
		fmt.Errorf("error execute template: %w", err)
		return
	}

	html := buf.String()

	m := gomail.NewMessage()
	for _, r := range receiverList {
		m.SetHeader("Subject", r.Subject)
		m.SetHeader("From", loadedConfig.EmailSenderName)
		m.SetAddressHeader("To", r.Address, r.Name)
		m.SetBody("text/html", html)

		if err := gomail.Send(s, m); err != nil {
			fmt.Printf("Could not send email to %q: %v", r.Address, err)
			return
		}
		m.Reset()
	}

	fmt.Println("Success sending email")
}
