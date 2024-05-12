ARG GO_VERSION=1
FROM golang:${GO_VERSION}-bookworm as builder

WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .

RUN go build ./src/cmd/main.go

FROM debian:bookworm

COPY --from=builder /usr/src/app/main /usr/local/bin/
COPY --from=builder /usr/src/app/src/migrations /usr/src/app/migrations

CMD ["main"]
