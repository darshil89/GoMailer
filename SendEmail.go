package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"

	"mailingService/models"

	"github.com/wneessen/go-mail"
)

var data models.EmailData

func sendEmail(w http.ResponseWriter, r *http.Request) {

	// extract data from request body
	err := json.NewDecoder(r.Body).Decode(&data)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid format")
		log.Fatalf("Failed to decode request body: %s", err)
		return
	}

	// print the data to the console
	log.Printf("Data: %v", data)

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
	// Parse the HTML template
	tpl, err := template.ParseFiles("email_template.html")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to parse HTML template")
		log.Fatalf("Failed to parse HTML template: %s", err)
		return
	}

	// Set the email body to the generated HTML content using the template and data
	if err := m.SetBodyHTMLTemplate(tpl, data); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to set HTML template mail body")
		log.Fatalf("Failed to set HTML template mail body: %s", err)
		return
	}

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
