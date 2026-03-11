## Stage 1
FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o /app/scheduler ./main.go

## Stage 2
FROM ubuntu:latest

WORKDIR /app

RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/scheduler /app/scheduler
COPY web /app/web

ENV TODO_PORT=":7540"
ENV TODO_DBFILE="/data/scheduler.db"

EXPOSE 7540

CMD ["/app/scheduler"]

