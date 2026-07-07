FROM golang:1.26.4-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o worker ./cmd/worker

EXPOSE 9091

CMD ["./worker"]