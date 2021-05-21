FROM golang:1.16.4-alpine3.13

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o app ./cmd/api/main.go

CMD ["./app"]
