# Estágio de compilação
FROM golang:1.17-alpine AS builder

# Instalar as dependências necessárias para o build do librdkafka e confluent-kafka-go
RUN apk add --no-cache \
    build-base \
    gcc \
    git \
    cmake \
    openssl-dev \
    bash

# Clonar e compilar librdkafka
WORKDIR /librdkafka
RUN git clone --branch v1.9.2 --depth 1 https://github.com/edenhill/librdkafka.git . && \
    ./configure --prefix /usr && make && make install && rm -rf /librdkafka


# FROM golang:1.17-alpine AS builder

# RUN apk add --no-cache \
#     bash \ 
#     build-base \
#     git \
#     cmake \
#     openssl-dev \
#     musl-dev

# WORKDIR /librdkafka
# RUN git clone --branch v1.9.2 --depth 1 https://github.com/edenhill/librdkafka.git . && \
#     ./configure --prefix /usr LDFLAGS=-static-libgcc && \
#     make && make install && \
#     rm -rf /librdkafka

WORKDIR /build

# Voltar para o diretório de build da aplicação Golang
WORKDIR /app

# Copiar e baixar as dependências da aplicação
COPY go.mod go.sum ./
RUN go mod download

# Copiar o código-fonte da aplicação
COPY . .

# COPY go.mod go.sum ./
# RUN go mod download

RUN go get -u github.com/confluentinc/confluent-kafka-go@latest

RUN CGO_ENABLED=1 GOOS=linux go build -tags musl -o app ./balance-api/cmd/consumer

# Estágio final
FROM alpine:latest

# Instalar a biblioteca librdkafka
RUN apk add --no-cache \
    librdkafka

# Copiar a aplicação compilada
COPY --from=builder /app /app

# Definir o diretório de trabalho
WORKDIR /

# Executar a aplicação
CMD ["/app"]
