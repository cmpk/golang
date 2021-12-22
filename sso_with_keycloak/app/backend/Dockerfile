FROM golang:1.17.4-bullseye

ENV GOPATH=

# 自己証明書対応
COPY ./ssl/server.crt /usr/local/share/ca-certificates/server.crt
RUN apt update && \
    apt install -y ca-certificates && \
    update-ca-certificates

### Docker build 時点では volume がマウントされていないため、必要なファイルがあればコピーする必要がある。
COPY go.mod go.sum ./
RUN go mod download

WORKDIR /go/app/cmd

### ホットリロードが不要な場合
# CMD ["go", "run", "main.go"]

# 開発環境でホットリロードさせる
RUN go get -u github.com/cosmtrek/air && \
    go build -o /go/bin/air github.com/cosmtrek/air

CMD ["air"]
