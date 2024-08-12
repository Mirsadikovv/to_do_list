FROM golang:1.22.1-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY init.sql /docker-entrypoint-initdb.d/

COPY . .

RUN go build -o go_user_service ./cmd/main.go

RUN chmod +x start.sh

EXPOSE 8081

CMD ["./start.sh"]
