#!/usr/bin/env bash

function build_sampleplugin() {
pushd "$THIS_PROJECT"/src/plugin-manager/fixture/sampleplugin/testplugin
  GOOS=linux go build -o helloworld-linux *.go
popd
}

function package_sampleplugin() {
pushd "$THIS_PROJECT"/src/plugin-manager/fixture/sampleplugin
  zip -r testplugin.zip testplugin
popd
}

build_sampleplugin
package_sampleplugin

