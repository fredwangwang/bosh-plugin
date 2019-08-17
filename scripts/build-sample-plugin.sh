#!/usr/bin/env bash

ROOT="$(git rev-parse --show-toplevel)"

function build_sampleplugin() {
pushd "$ROOT"/example/sample-plugin
  GOOS=linux go build -o main-linux *.go
popd
}

function package_sampleplugin() {
pushd "$ROOT"/example
  rm sample-plugin.zip &> /dev/null || true
  zip -r -j sample-plugin.zip sample-plugin/*
popd
}

build_sampleplugin
package_sampleplugin

