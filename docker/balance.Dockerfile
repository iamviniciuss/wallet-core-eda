FROM golang:1.20 AS builder

WORKDIR /build

RUN apt-get update && apt-get install -y librdkafka-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

COPY docker/entrypoint-ba.sh /entrypoint-ba.sh
RUN chmod +x /entrypoint-ba.sh

RUN CGO_ENABLED=1 GOOS=linux go build -tags dynamic -o balanceapi ./balance-api/cmd/consumer
RUN CGO_ENABLED=1 GOOS=linux go build -o balanceapi_scripts ./balance-api/scripts/database

EXPOSE 3003

RUN chmod +x /entrypoint-ba.sh
ENTRYPOINT ["/entrypoint-ba.sh"]

