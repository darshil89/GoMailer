ARG GO_VERSION=1.23.5
FROM golang:${GO_VERSION}-alpine

# Set the working directory
WORKDIR /src

# Install runtime dependencies
RUN apk --update add ca-certificates tzdata && update-ca-certificates

# Copy Go modules and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the application source code
COPY . .

# Copy the email template
COPY email_template.html /src/email_template.html

# Build the application
ARG TARGETARCH
RUN CGO_ENABLED=0 GOARCH=$TARGETARCH go build -o /bin/server .

# Ensure the binary is executable
RUN chmod +x /bin/server

# Expose the application port
EXPOSE 8080

# Run the application
ENTRYPOINT ["/bin/server"]
