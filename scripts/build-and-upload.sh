#!/usr/bin/env bash

ROOT="$(git rev-parse --show-toplevel)"

export GOOS=linux

pushd "$ROOT/src/plugin-manager"
  go build -o plugin-manager-linux
  zip plugin.zip plugin-manager-linux

  bosh -d bosh-plugin scp plugin.zip bosh-plugin/0:/tmp/plugin.zip
  rm plugin.zip plugin-manager-linux
popd


# remote:
# monit stop plugin-manager && unzip -o /tmp/plugin.zip && cp plugin-manager-linux /var/vcap/packages/plugin-manager/plugin-manager && monit start plugin-manager
