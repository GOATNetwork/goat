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

That gives you a web3 JSON-RPC on port 8545 and a consensus REST API on port 1317.

### Running a local build against regtest

`make init` builds `goatd` from the pinned submodule, so to exercise a working
tree instead you have to drive the same steps by hand. Everything below runs
from the `goat-regtest` checkout; point `submodule/goat` at your tree, or just
drop your own binary into `build/`.

Prerequisites: `go`, `node`, `npm`, `jq`, and `docker` — the contract genesis
task starts an anvil container.

**1. Build the binaries.** `geth` comes from
[goat-geth](https://github.com/GOATNetwork/goat-geth); `goatd` is this repo.

```sh
mkdir -p build data/goat data/geth
make -C submodule/geth geth && cp submodule/geth/build/bin/geth build/
make -C submodule/goat build && cp submodule/goat/build/goatd build/
```

**2. Create the node home and the validator and voter keys.** `modgen init`
writes the config and the privval key; `locking sign` produces the proof of
possession that `Locking.sol` checks at registration.

```sh
cp example.json config.json
./build/goatd --home ./data/goat modgen init --regtest regtest

OWNER=0xbc000FE892bC88F2ba41d70aF9F80619F556dCA2
VALIDATOR=$(./build/goatd --home ./data/goat modgen locking sign --owner $OWNER)
jq --argjson v "$VALIDATOR" '.Locking.validators += [$v]' config.json > tmp.json && mv tmp.json config.json

VOTER=$(./build/goatd --home ./data/goat modgen relayer keygen --output 1.json)
jq --argjson v "$VOTER" '.Relayer.voters += [$v]' config.json > tmp.json && mv tmp.json config.json
```

**3. Generate the execution layer genesis and initialize geth.** Pass an
absolute path to `--param`: the task resolves it relative to the contracts
directory, which breaks if `submodule/contracts` is a symlink.

```sh
npm ci --prefix submodule/contracts --engine-strict
npm run compile --prefix submodule/contracts
npm run genesis --prefix submodule/contracts -- --param "$PWD/config.json"

./build/geth init --state.scheme hash --cache.preimages \
  --datadir ./data/geth ./submodule/contracts/genesis/regtest.json
```

**4. Assemble the consensus genesis.** This reads `config.json` and calls
`modgen` once per token, validator, voter and parameter group. It ends by
printing the execution layer genesis hash, which must match the one `geth
init` reported.

```sh
./submodule/goat/contrib/scripts/genesis.sh ./data/goat ./config.json
```

**5. Start both processes.** They talk over the geth IPC socket, so geth has to
be up first.

```sh
./build/geth --datadir ./data/geth --gcmode=archive --goat.preset=rpc --nodiscover &
./build/goatd start --home ./data/goat --regtest --goat.geth ./data/geth/geth.ipc &
```

Blocks start within a few seconds. The two heights advance in lockstep,
so comparing them is the quickest check that the engine API handshake is
healthy:

```sh
# consensus height
curl -s localhost:26657/status | jq -r '.result.sync_info.latest_block_height'

# execution height, same value in hex
curl -s -X POST -H 'Content-Type: application/json' \
  --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' \
  localhost:8545 | jq -r '.result'

# the validator, whose reward should be growing every block
curl -s "localhost:1317/goatnetwork/goat/locking/v1/validator?address=$(jq -r '.Locking.validators[0].validator' config.json)" | jq .
```

Note that a single validator network never produces a `Nil` precommit, so
anything that only shows up when the votes disagree needs more than one node.
