FROM golang:1.23.3-bookworm
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 go build -o bin/server cmd/server/main.go
CMD ./bin/server