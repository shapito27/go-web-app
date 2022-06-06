package main

import (
	"time"

	"github.com/shapito27/go-web-app/internal/models"
	mail "github.com/xhit/go-simple-mail/v2"
)

func listenForMail() {
	go func() {
		for {
			mailData := <-app.MailChan
			sendMsg(mailData)
		}
	}()
}

func sendMsg(m models.MailData) {
	server := mail.NewSMTPClient()
	server.Host = "localhost"
	server.Port = 1025
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	client, err := server.Connect()
	if err != nil {
		app.ErrorLog.Println(err)
	}

	email := mail.NewMSG()
	email.SetSubject(m.Subject)
	email.SetFrom(m.From)
	email.AddTo(m.To)
	email.SetBody(mail.TextHTML, m.Content)
	err = email.Send(client)
	if err != nil {
		app.ErrorLog.Println(err)
	}

}
