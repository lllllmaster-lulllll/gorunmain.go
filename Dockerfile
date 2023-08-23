From golang:1.21.0 AS builder

WORKDIR /app

ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY ./main.go .
COPY ./myapp/ ./myapp

RUN go mod tidy

RUN go build -o /myapp

ENTRYPOINT [ "/myapp" ]