package main

import (
	"log"
	"net/http"
	"os"

	"mailingService/models"

	"github.com/wneessen/go-mail"
)

var data models.EmailData

func sendEmail(w http.ResponseWriter, r *http.Request) {

	//Setting up the message body
	var email = os.Getenv("EMAIL")
	var password = os.Getenv("PASSWORD")
	m := mail.NewMsg()
	if err := m.From(email); err != nil {
		log.Fatalf("Failed to set FROM address %s", err)
		return
	}
	if err := m.To(data.Email); err != nil {
		log.Fatalf("Failed to set TO address %s", err)
		return
	}
	m.Subject("Thankyou for contacting me")
	reply := "Hi " + data.Name + ",\n I read your message and will reach you back ASAP! :-D"
	m.SetBodyString(mail.TypeTextPlain, reply)

	//Setting up mail client

	c, err := mail.NewClient("smtp.gmail.com", mail.WithPort(587), mail.WithSMTPAuth(mail.SMTPAuthPlain), mail.WithTLSPortPolicy(mail.TLSMandatory), mail.WithUsername(email), mail.WithPassword(password))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create client")
		log.Fatalf("Failed to create client %s", err)
		return
	}
	//Sending the email
	if err := c.DialAndSend(m); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to send email")
		log.Fatalf("Failed to send mail: %s", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, "Message Created")
}
