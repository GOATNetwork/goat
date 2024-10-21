#!/bin/bash

# the genesis config file 
# https://github.com/GOATNetwork/goat-contracts/blob/main/task/deploy/param.ts

set -e

if [ ! -d "$1" ]; then
  echo "goat home $1 is not a directory"
  exit 1
fi

if [ ! -f "$2" ]; then
  echo "genesis config $2 is not a file"
  exit 1
fi

if [ ! -f "$3" ]; then
  echo "geth-genesis $3 is not a file"
  exit 1
fi

TOKEN_LENGTH=$(jq '.Locking.tokens | length' $2)
for ((i=0; i<TOKEN_LENGTH; i++)); do
  echo "Add token $i"
  address=$(jq -r ".Locking.tokens[$i].address" $2)
  echo "address: $address"
  weight=$(jq -r ".Locking.tokens[$i].weight" $2)
  echo "weight: $weight"
  threshold=$(jq -r ".Locking.tokens[$i].threshold" $2)
  echo "threshold: $threshold"

  ./build/goatd modgen locking add-token --home $1 --token $address --weight $weight --threshold $threshold
done

VALIDATOR_LENGTH=$(jq '.Locking.validators | length' $2)
for ((i=0; i<VALIDATOR_LENGTH; i++)); do
  echo "Add validator $i"
  pubkey=$(jq -r ".Locking.validators[$i].pubkey" $2)
  echo "pubkey: $pubkey"
  ./build/goatd modgen locking add-validator --home $1 --pubkey $pubkey
done

GOAT_LOCKING_CONTRACT='0xbc10000000000000000000000000000000000004'
TRANSFER_LENGTH=$(jq '.GoatToken.transfers | length' $2)
for ((i=0; i<TRANSFER_LENGTH; i++)); do
  address=$(jq -r ".GoatToken.transfers[$i].to" $2 | tr '[:upper:]' '[:lower:]')
  if [ $address = $GOAT_LOCKING_CONTRACT ]; then
    value=$(jq -r ".GoatToken.transfers[$i].value" $2)
    echo "Set initial reward $value"
    ./build/goatd modgen locking --home $1 --init-reward $value
  fi
done

./build/goatd modgen locking param --home $1 \
  --unlock-duration $(jq -r ".Consensus.Locking.unlockDuration" $2) \
  --exit-duration $(jq -r ".Consensus.Locking.exitDuration" $2)

VOTERS_LENGTH=$(jq '.Relayer.voters | length' $2)
for ((i=0; i<VOTERS_LENGTH; i++)); do
  echo "Add voter $i"
  txKey=$(jq -r ".Relayer.voters[$i].txKey" $2)
  echo "txKey: $txKey"
  voteKey=$(jq -r ".Relayer.voters[$i].voteKey" $2)
  echo "voteKey: $voteKey"
  ./build/goatd modgen relayer add-voter --home $1 --key.tx $txKey --key.vote $voteKey
done

DEPOSIT_LENGTH=$(jq '.Bridge.deposits | length' $2)
for ((i=0; i<DEPOSIT_LENGTH; i++)); do
  echo "Add deposit $i"
  txid=$(jq -r ".Bridge.deposits[$i].txid" $2)
  echo "txid: $txid"
  txout=$(jq -r ".Bridge.deposits[$i].txout" $2)
  echo "txout: $txout"
  satoshi=$(jq -r ".Bridge.deposits[$i].satoshi" $2)
  echo "satoshi: $satoshi"
  ./build/goatd modgen bitcoin add-deposit --home $1 --txid $txid --txout $txout --satoshi $satoshi
done

./build/goatd modgen relayer --home $1 --param.accept_proposer_timeout $(jq -r ".Consensus.Relayer.acceptProposerTimeout" $2)

./build/goatd modgen bitcoin --home $1 \
  --min-deposit $(jq -r ".Consensus.Bridge.minDepositInSat" $2) \
  --confirmation-number $(jq -r ".Consensus.Bridge.confirmationNumber" $2) \
  --network $(jq -r ".Bitcoin.network" $2) \
  --pubkey $(jq -r ".Consensus.Relayer.tssPubkey" $2) \
  $(jq -r ".Bitcoin.height" $2) $(jq -r ".Bitcoin.hash" $2)

./build/goatd modgen goat --home $1 $3
