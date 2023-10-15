FROM golang:1.21.3-alpine3.17 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY main.go/ ./
COPY game/ ./game
COPY magic/ ./magic
COPY server/ ./server
COPY models/ ./models
COPY stockfish/ ./stockfish

RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/server

FROM alpine:3.17

WORKDIR /bin

COPY .ssh/ /bin
COPY stockfish/ /bin/stockfish
COPY --from=build /bin/server /bin/server

EXPOSE 2324

CMD ["/bin/server"]
