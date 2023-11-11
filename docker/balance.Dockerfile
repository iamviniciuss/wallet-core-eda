FROM golang:1.20 AS builder

WORKDIR /build

RUN apt-get update && apt-get install -y librdkafka-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -tags dynamic -o balanceapi ./balance-api/cmd/consumer

EXPOSE 3003

ENTRYPOINT ["./balanceapi"]

