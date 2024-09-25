FROM golang:latest

WORKDIR /app

COPY go.mod go.sum .
RUN go mod download

# ensure source directories are copied to target directories individually
COPY cmd/ cmd/
COPY pkg/ pkg/
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/server ./cmd/

CMD ["/app/bin/server"]
