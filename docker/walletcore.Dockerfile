FROM golang:1.20 AS builder

WORKDIR /build

RUN apt-get update && apt-get install -y librdkafka-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

COPY docker/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

RUN go build -tags dynamic -o walletcore_app ./wallet-core/cmd/walletcore
RUN go build -o walletcore_scripts ./wallet-core/scripts/database

EXPOSE 8080

# ENTRYPOINT ["./walletcore_app"]
RUN chmod +x /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]

