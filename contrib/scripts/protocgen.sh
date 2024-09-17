#!/usr/bin/env bash

set -ex

echo "Formatting protobuf files"
find ./ -name "*.proto" -exec clang-format -i {} \;
buf format -w

home=$PWD

echo "Generating proto code"
proto_dirs=$(find . -name 'buf.yaml' -print0 | xargs -0 -n1 dirname | sort | uniq)
for dir in $proto_dirs; do
  echo "Generating proto code for $dir"

  cd $dir
  # check if buf.gen.pulsar.yaml exists in the proto directory
  if [ -f "buf.gen.pulsar.yaml" ]; then
    buf generate --template buf.gen.pulsar.yaml
  fi

  # check if buf.gen.gogo.yaml exists in the proto directory
  if [ -f "buf.gen.gogo.yaml" ]; then
      for file in $(find . -maxdepth 8 -name '*.proto'); do
        if grep -q "option go_package" "$file"; then
          buf generate --template buf.gen.gogo.yaml $file
        fi
    done

    # move generated files to the right places
    if [ -d "../github.com" ]; then
      cp -r ../github.com/goatnetwork/goat/* $home
      rm -rf ../github.com
    fi
  fi

  cd $home
done

# move generated files to the right places
rm -rf github.com
