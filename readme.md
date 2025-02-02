# 📧 GoMailService - Go-based Mailing Microservice

Welcome to **GoMailService**! This is a simple mailing microservice built using Go, designed to send emails and handle basic email validation. This project demonstrates how to create a REST API in Go that can handle email requests and verify the provided email addresses.

---

## 🚀 Features
- 🖥 **Go-based HTTP Server** for handling email requests
- 📧 **Send Emails** through a POST endpoint
- 🌐 **Email Validation** before sending emails
- 🔒 **Basic Authentication** for added security
- 🔍 **Email Verification** using an external service
- 🛠 **Error Handling** for invalid email formats and failures during email sending

---

## 🛠️ Getting Started

### Prerequisites

Before getting started, make sure you have Go installed. You can download and install Go from the official Go website:

- [Download Go](https://golang.org/dl/)

## 🌍 API Endpoints

### 1. **POST /api/getToken**

This endpoint creates a token for authentication. Use this token to authenticate requests to the `sendEmail` endpoint.

**Request Body:**
```json
{
  "secret": "your-username"
}
```

### 2. **POST /api/sendEmail**

This endpoint send mail using the secret token

**Request Body:**
```json
{
  "email": "recipient@example.com",
  "subject": "Subject of the Email",
  "message": "Hello, this is the email content!"
}
```

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

Thank you for checking out **GoMailService**! 🎉