FROM golang:1.23-alpine

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o api-gateway cmd/main.go

EXPOSE 5000

CMD ["./api-gateway"]
