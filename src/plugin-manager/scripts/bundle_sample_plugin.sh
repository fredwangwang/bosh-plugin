#!/usr/bin/env bash

function build_sampleplugin() {
pushd "$THIS_PROJECT"/src/plugin-manager/fixture/sampleplugin/testplugin
  GOOS=linux go build -o helloworld-linux *.go
popd
}

function package_sampleplugin() {
pushd "$THIS_PROJECT"/src/plugin-manager/fixture/sampleplugin
  rm testplugin.zip
  zip -r -j testplugin.zip testplugin/*
popd
}

build_sampleplugin
package_sampleplugin

