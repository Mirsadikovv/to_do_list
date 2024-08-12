FROM golang:1.22.1-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o api_gateway ./cmd/main.go

EXPOSE 8080

CMD ["./api_gateway"]
