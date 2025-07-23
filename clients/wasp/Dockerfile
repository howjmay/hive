# Stage 1: Build evmemulator from source
FROM golang:1.23 AS evmbuilder

WORKDIR /src
RUN git clone --depth=1 --branch emulator https://github.com/howjmay/wasp.git waspclone
# COPY waspclone /src/wasp
WORKDIR /src/waspclone/tools/evm/evmemulator
RUN go build -o /evmemulator

# Stage 2: Final runtime image
FROM ubuntu:22.04

RUN apt-get update && apt-get install -y curl tar ca-certificates libpq5 \
    && rm -rf /var/lib/apt/lists/*

# Download and extract IOTA release
ENV IOTA_VERSION=v1.2.3
RUN curl -L \
  https://github.com/iotaledger/iota/releases/download/${IOTA_VERSION}/iota-${IOTA_VERSION}-linux-x86_64.tgz \
  -o /tmp/iota.tgz \
 && tar -xzf /tmp/iota.tgz -C /usr/local/bin \
 && rm /tmp/iota.tgz

# Copy evmemulator binary from builder stage
COPY --from=evmbuilder /evmemulator /usr/local/bin/evmemulator

# Genesis file
COPY genesis.json /etc/genesis.json

# Entrypoint script
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

EXPOSE 9000 9123 8545

ENTRYPOINT ["/entrypoint.sh"]