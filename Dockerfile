FROM ubuntu:24.04

RUN apt update && apt install -y \
    curl \
    vim \
    wget \
    && apt update \
    && rm -rf /var/lib/apt/lists/*

# Download Golang
RUN wget https://go.dev/dl/go1.24.2.linux-amd64.tar.gz \
    && tar -C /usr/local -xzf go1.24.2.linux-amd64.tar.gz \
    && rm -rf go1.24.2.linux-amd64.tar.gz

ENV PATH="/usr/local/go/bin:${PATH}"

# 作業ディレクトリ作成
WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .

RUN go build -o app .

EXPOSE 7777

# アプリケーションを実行
CMD ["./app"]