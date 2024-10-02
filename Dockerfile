FROM golang:latest

# Install air for live reloading
RUN go install github.com/air-verse/air@latest

WORKDIR /app

COPY go.mod go.sum .
RUN go mod download

EXPOSE 2024

CMD ["air", "-c", ".air.toml"]
