#!/usr/bin/env bash

ROOT="$(git rev-parse --show-toplevel)"

export GOOS=linux

pushd "$ROOT/src/plugin-manager"
  go build -o plugin-manager-linux
  zip plugin.zip plugin-manager-linux

  bosh -d bosh-plugin scp plugin.zip bosh-plugin/0:/tmp/plugin.zip
#  bosh -d bosh-plugin ssh bosh-plugin/0 --non-interactive --command  \
#      'sudo -i && cp /tmp/plugin-manager-linux /var/vcap/packages/plugin-manager/plugin-manager && exit'
popd


# remote:
# monit stop plugin-manager && unzip -o /tmp/plugin.zip && cp /tmp/plugin-manager-linux /var/vcap/packages/plugin-manager/plugin-manager && monit start plugin-manager

 METRON_ADDR, METRON_CA, METRON_CERT, METRON_KEY, MONIT