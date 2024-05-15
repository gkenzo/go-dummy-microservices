FROM golang:1.22.2-alpine as builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o frontApp ./cmd/web

RUN chmod +x /app/frontApp

FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/frontApp /app

COPY ./cmd/web/templates /app/cmd/web/templates

CMD ["/app/frontApp"]
