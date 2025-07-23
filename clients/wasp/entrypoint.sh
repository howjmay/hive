#!/bin/bash
set -e

# Start IOTA node with faucet
/usr/local/bin/iota start \
  --force-regenesis \
  --epoch-duration-ms=60000 \
  --with-faucet \
  --faucet-amount=1000000000000 &
IOTA_PID=$!

# Wait until IOTA is responsive
echo "Waiting for IOTA on port 9000..."
until curl -sf \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","id":1,"method":"iota_getLatestCheckpointSequenceNumber","params":[]}' \
  http://localhost:9000; do
  sleep 2
done
echo "✔ IOTA is ready!"

# Launch evmemulator in foreground
exec evmemulator --node-launch-mode=docker-compose --genesis=/etc/genesis.json
