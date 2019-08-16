#!/usr/bin/env bash

ROOT="$(git rev-parse --show-toplevel)"

set -eu

pushd "$ROOT"
bosh deploy --deployment bosh-plugin \
     manifest/manifest.yml \
     --vars-file <(lpass show --notes plugin-manager.yml)
popd