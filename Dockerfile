FROM golang:latest as builder

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

RUN go build -o go-tcp-chat-server .

FROM gcr.io/distroless/base

WORKDIR /

COPY --from=builder /app/go-tcp-chat-server .

CMD ["./go-tcp-chat-server"]