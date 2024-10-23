# Description: Dockerfile for building a go binary and running it in a scratch container

FROM golang:1.23.2-alpine3.9 as builder

WORKDIR /app

COPY go.mod .

COPY go.sum .

# go mod tidy will remove any dependencies that are no longer needed
RUN go mod tidy 

# go mod download will download all dependencies
RUN go mod download

COPY . .

RUN GOOS=linux GOARCH=amd64 go build -o main .

FROM scratch

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 3000

# Run the binary
CMD ["./main"]