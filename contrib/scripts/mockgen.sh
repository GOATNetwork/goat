#!/usr/bin/env bash

mockgen_cmd="mockgen"
$mockgen_cmd -source=x/relayer/types/types.go -destination=testutil/mock/relayer_vote_msg.go -package=mock IVoteMsg
$mockgen_cmd -source=pkg/ethrpc/interface.go -destination=testutil/mock/eth_engine_client.go -package=mock EngineClient
$mockgen_cmd -source=x/goat/types/expected_keepers.go -destination=testutil/mock/locking_keeper.go -package=mock -exclude_interfaces RelayerKeeper,BitcoinKeeper,AccountKeeper LockingKeeper
$mockgen_cmd -source=x/goat/types/expected_keepers.go -destination=testutil/mock/relayer_keeper.go -package=mock -exclude_interfaces LockingKeeper,BitcoinKeeper,AccountKeeper RelayerKeeper
$mockgen_cmd -source=x/goat/types/expected_keepers.go -destination=testutil/mock/bitcoin_keeper.go -package=mock -exclude_interfaces LockingKeeper,RelayerKeeper,AccountKeeper BitcoinKeeper
$mockgen_cmd -source=x/goat/types/expected_keepers.go -destination=testutil/mock/account_keeper.go -package=mock -exclude_interfaces LockingKeeper,RelayerKeeper,BitcoinKeeper AccountKeeper
