FROM golang:1.22.0-alpine as builder

LABEL maintainer="Sultonov Daler <daler.sultonov.00@gmail.com>"

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o main cmd/main.go


FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/.env.json .

CMD ["./main"]