# Estágio de compilação
FROM golang:1.20 AS builder

WORKDIR /build

# Instale as dependências necessárias e remova os pacotes não necessários
RUN apt-get update \
    && apt-get install -y librdkafka-dev \
    && go mod download \
    && apt-get purge -y --auto-remove \
    && rm -rf /var/lib/apt/lists/*

# Copie apenas os arquivos necessários para baixar as dependências e compilar
COPY go.mod go.sum ./
COPY . .

# Copie o script de entrypoint e ajuste permissões
COPY docker/entrypoint-ba.sh /entrypoint-ba.sh
RUN chmod +x /entrypoint-ba.sh

# Compile a aplicação
RUN CGO_ENABLED=1 GOOS=linux go build -tags dynamic -o balanceapi ./balance-api/cmd/consumer \
    && CGO_ENABLED=1 GOOS=linux go build -o balanceapi_scripts ./balance-api/scripts/database

# Estágio final
FROM alpine:latest

WORKDIR /app

# Copie apenas os artefatos necessários do estágio de compilação
COPY --from=builder /build/balanceapi .
COPY --from=builder /build/balanceapi_scripts .

EXPOSE 3003

# Copie o script de entrypoint e ajuste permissões
COPY --from=builder /entrypoint-ba.sh /entrypoint-ba.sh
RUN chmod +x /entrypoint-ba.sh

ENTRYPOINT ["/entrypoint-ba.sh"]
