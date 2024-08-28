FROM golang:1.21.1-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o go-echo cmd/server/main.go

EXPOSE 8080

CMD ["./go-echo"]