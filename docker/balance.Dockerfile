FROM alpine:latest AS builder
ENV LIBRD_VER=1.9.2
WORKDIR /tmp

RUN apk add --no-cache --virtual .make-deps bash make wget git gcc g++ && apk add --no-cache musl-dev zlib-dev openssl zstd-dev pkgconfig libc-dev && wget https://github.com/edenhill/librdkafka/archive/v${LIBRD_VER}.tar.gz && tar -xvf v${LIBRD_VER}.tar.gz && cd librdkafka-${LIBRD_VER} && ./configure --prefix /usr && make && make install && make clean && rm -rf librdkafka-${LIBRD_VER} && rm -rf v${LIBRD_VER}.tar.gz && apk del .make-deps

RUN export PKG_CONFIG_PATH=$PKG_CONFIG_PATH:/usr/lib/pkgconfig/

ENV PKG_CONFIG_PATH=$PKG_CONFIG_PATH:/usr/lib/pkgconfig/

RUN apk --no-cache add \
    bash \
    ca-certificates \
    gcc \
    musl-dev \
    librdkafka-dev 

RUN apk --no-cache add go

WORKDIR $GOPATH

RUN go version

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download
COPY . .

COPY docker/entrypoint-ba.sh ./entrypoint-ba.sh
RUN chmod +x ./entrypoint-ba.sh

RUN go build -tags dynamic -o balanceapi ./balance-api/cmd/consumer
RUN go build -o balanceapi_scripts ./balance-api/scripts/database




FROM alpine:latest 
ENV LIBRD_VER=1.9.2
WORKDIR /tmp
RUN apk add --no-cache --virtual .make-deps bash make wget git gcc g++ && apk add --no-cache musl-dev zlib-dev openssl zstd-dev pkgconfig libc-dev && wget https://github.com/edenhill/librdkafka/archive/v${LIBRD_VER}.tar.gz && tar -xvf v${LIBRD_VER}.tar.gz && cd librdkafka-${LIBRD_VER} && ./configure --prefix /usr && make && make install && make clean && rm -rf librdkafka-${LIBRD_VER} && rm -rf v${LIBRD_VER}.tar.gz && apk del .make-deps
RUN export PKG_CONFIG_PATH=$PKG_CONFIG_PATH:/usr/lib/pkgconfig/
ENV PKG_CONFIG_PATH=$PKG_CONFIG_PATH:/usr/lib/pkgconfig/
RUN apk --no-cache add bash 

WORKDIR /workspace

COPY --from=builder /build/entrypoint-ba.sh /workspace/entrypoint-ba.sh
COPY --from=builder /build/balanceapi /workspace/balanceapi
COPY --from=builder /build/balanceapi_scripts /workspace/balanceapi_scripts

RUN chmod +x /workspace/entrypoint-ba.sh
RUN chmod +x /workspace/balanceapi
RUN chmod +x /workspace/balanceapi_scripts

ENTRYPOINT ["/workspace/entrypoint-ba.sh"]
# ENTRYPOINT ["/bin/sh", "/app/entrypoint-ba.sh"]

