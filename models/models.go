package models

type Credentials struct {
	Secret string `json:"secret"`
}

type EmailData struct {
	Email   string `json:"email"`
	Message string `json:"message"`
	Name    string `json:"name"`
}
