FROM golang:1.21-alpine as builder

RUN export GOPROXY=direct

WORKDIR $GOPATH/src/app

RUN apk update && \
    apk add pkgconfig && \
    apk add vips-dev && \
    apk add build-base

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /server ./server/main.go


FROM alpine:latest

WORKDIR /app

COPY --from=builder /server .