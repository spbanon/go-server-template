FROM golang:1.21-alpine as builder

RUN export GOPROXY=direct

WORKDIR $GOPATH/src/app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /dbhandler ./dbhandler/dbhandler.go


FROM alpine:latest

WORKDIR /app
COPY --from=builder /dbhandler .