FROM golang:latest

WORKDIR /app


ENV GO111MODULE=on \
    GOOS=linux \
    GOARCH=amd64 \
    PORT=8080 \
    CGO_ENABLED=0

COPY . .

RUN go mod init main
RUN go mod tidy
RUN go build




CMD [“./gopherator“]
