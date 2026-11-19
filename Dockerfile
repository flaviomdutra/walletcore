FROM golang:1.27

WORKDIR /app/

RUN apt-get update && apt-get install -y librdkafka-dev

COPY go.mod go.sum ./
RUN go mod download

CMD ["go", "run", "cmd/walletcore/main.go"]
