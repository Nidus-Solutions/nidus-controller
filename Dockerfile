FROM golang:1.23.2-alpine3.9 as builder

WORKDIR /app

COPY go.mod .

COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o main .

FROM alpine:3.9

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 3000

CMD ["./main"]