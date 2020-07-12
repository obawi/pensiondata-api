FROM golang:1.14-alpine AS builder
ENV GO111MODULE=on
RUN apk add --no-cache git

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN GOOS=linux go build -o pensiondata-api *.go

FROM alpine:latest
COPY --from=builder /app/pensiondata-api /app/
CMD ["/app/pensiondata-api"]