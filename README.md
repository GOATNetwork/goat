# GOAT

GOAT consensus client built using Cosmos SDK

## Public networks

To participate in or utilize the public networks, take a look at the [the docker-compose files](./docker).

## Regtest

https://github.com/GOATNetwork/goat-regtest

```sh
git clone --recurse-submodules https://github.com/GOATNetwork/goat-regtest.git
cd goat-regtest
make init start
```

Web3 JSON-RPC on 8545, consensus REST on 1317.

### Producing blocks from a working tree

`make init` builds `goatd` from the pinned submodule, and `make start` runs it
under pm2. To test your own checkout instead, swap the binary in and run the
two processes directly.

Needs `go`, `node`, `npm`, `jq` and `docker` — the contract genesis task starts
an anvil container.

Run from the root of this repo:

```sh
GOAT=$PWD
make build
git clone --recurse-submodules https://github.com/GOATNetwork/goat-regtest.git ../goat-regtest
cd ../goat-regtest
mkdir -p build data/goat data/geth
make geth contracts
cp $GOAT/build/goatd build/
cp example.json config.json
sh ./init.sh
./build/geth --datadir ./data/geth --gcmode=archive --goat.preset=rpc --nodiscover &
while [ ! -S data/geth/geth.ipc ]; do sleep 1; done
./build/goatd start --home ./data/goat --regtest --goat.geth ./data/geth/geth.ipc &
```

Blocks start within seconds. Check that consensus and execution advance
together, which is the quickest sign the engine API handshake is healthy:

```sh
curl -s localhost:26657/status | jq -r .result.sync_info.latest_block_height
curl -s -X POST -H 'Content-Type: application/json' \
  --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' \
  localhost:8545 | jq -r .result
curl -s "localhost:1317/goatnetwork/goat/locking/v1/validator?address=$(jq -r .Locking.validators[0].validator config.json)" | jq .validator
```

Confirm you are running what you think you are:

```sh
./build/goatd version --long | grep -E '^version:|^cosmos_sdk_version:|^go:'
```

Stop `goatd` first, then give `geth` a few seconds to flush:

```sh
kill $(pgrep -x goatd) ; sleep 4 ; kill $(pgrep -x geth)
```

Never `kill -9` these. `goatd` commits a height as soon as `geth` accepts the
payload, so killing `geth` before it persists leaves the consensus head
pointing at an EVM block that no longer exists on disk. The chain then stops
at the next height with

```
ERR failed to prepare proposal err="failed to build goat-geth txs"
WARN Fetching the unknown forkchoice head from network
```

and the only way out is to start over:

```sh
rm -rf data/goat data/geth && mkdir -p data/goat data/geth && sh ./init.sh
```

Two things this setup cannot show you. A single validator never casts a `Nil`
precommit, so anything that only appears when votes disagree needs more nodes.
And a chain started from scratch never exercises an upgrade: to check that a
new binary can take over an existing chain, stop cleanly, replace
`build/goatd`, and start again over the same data.
