ARG GO_VERSION=1.23.5
FROM golang:${GO_VERSION}-alpine AS build

WORKDIR /src

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .
COPY email_template.html /src/email_template.html

# Build the application
ARG TARGETARCH
RUN CGO_ENABLED=0 GOARCH=$TARGETARCH go build -o /bin/server .

# Final runtime image
FROM alpine:latest

# Install runtime dependencies
RUN apk --update add ca-certificates tzdata && update-ca-certificates

# Copy the built binary from the build stage
COPY --from=build /bin/server /bin/

# Expose the port
EXPOSE 8080

# Run the application
ENTRYPOINT [ "/bin/server" ]
