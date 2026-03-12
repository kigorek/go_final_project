## Stage 1
FROM golang:1.25.1 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/scheduler ./main.go

## Stage 2
FROM alpine:latest

WORKDIR /app

ARG TODO_PORT=7540

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/scheduler /app/scheduler
COPY web /app/web

ENV TODO_PORT=":${TODO_PORT}"
ENV TODO_DBFILE="/data/scheduler.db"

EXPOSE ${TODO_PORT}

CMD ["/app/scheduler"]

