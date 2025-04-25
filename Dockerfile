FROM ubuntu:24.04

RUN apt update && apt install -y \
    curl \
    vim \
    wget \
    && curl -sSL https://ngrok-agent.s3.amazonaws.com/ngrok.asc \
    | tee /etc/apt/trusted.gpg.d/ngrok.asc >/dev/null \
    && echo "deb https://ngrok-agent.s3.amazonaws.com buster main" \
    | tee /etc/apt/sources.list.d/ngrok.list \
    && apt update \
    && apt install -y ngrok \
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

RUN go build -o app

# start.sh スクリプトをコピー
COPY start.sh /start.sh

# スクリプトを実行可能にする
RUN chmod +x /start.sh

EXPOSE 7777

CMD ["/start.sh"]