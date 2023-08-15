From golang:1.21.0 AS builder

WORKDIR /app

COPY go.mod go.sum ./
ENV GOPROXY https://goproxy.cn
ENV CGO_ENABLED 0 
ENV GOARCH amd64
ENV GOOS linux

RUN go mod download && go mod verify

COPY main.go ./

RUN go build -o /my-app

ENTRYPOINT [ "/my-app" ]