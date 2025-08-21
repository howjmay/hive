#!/bin/bash
set -e

# Start IOTA node with faucet
# /usr/local/bin/iota start \
#   --force-regenesis \
#   --epoch-duration-ms=60000 \
#   --with-faucet \
#   --faucet-amount=1000000000000 > /dev/null 2>&1 &
# IOTA_PID=$!

# # Wait until IOTA is responsive
# echo "Waiting for IOTA on port 9000..."
# until curl -sf \
#   -H "Content-Type: application/json" \
#   -d '{"jsonrpc":"2.0","id":1,"method":"iota_getLatestCheckpointSequenceNumber","params":[]}' \
#   http://localhost:9000; do
#   sleep 2
# done
# echo "✔ IOTA is ready!"

# until curl -sf \
#   -H "Content-Type: application/json" \
#   -d '{"FixedAmountRequest":{"recipient":"0x1a15c995b94cdff9291d9c7218be747d0b44b70f5849103c2683d540dd3f920b"}}' \
#   http://localhost:9123/gas; do
#   sleep 2
# done
# echo "✔ faucet is ready!"

# Launch evmemulator in foreground
# exec evmemulator --node-launch-mode=docker-compose --genesis=/etc/genesis.json
# '/genesis.json' is the path of genesis file that hive pass the from fixtures
exec evmemulator --node-launch-mode=docker-compose --remote-host=http://host.docker.internal --genesis=/genesis.json
