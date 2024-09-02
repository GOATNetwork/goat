FROM golang:1.22-alpine AS builder
RUN apk add --no-cache curl make git libc-dev bash file gcc linux-headers eudev-dev
WORKDIR /app
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
RUN LEDGER_ENABLED=false LINK_STATICALLY=true BUILD_TAGS=muslc make build
RUN echo "Ensuring binary is statically linked ..."  \
    && file /app/build/goatd | grep "statically linked"

FROM alpine:3.20
RUN apk add --no-cache build-base jq
RUN addgroup -g 1025 nonroot && \
    adduser -D nonroot -u 1025 -G nonroot
COPY --from=builder /app/build/goatd /usr/local/bin/
EXPOSE 26656 26657 1317 9090
ENTRYPOINT ["goatd", "start"]