#!/usr/bin/env bash

ROOT="$(git rev-parse --show-toplevel)"

pushd "$ROOT/example/sample-plugin-java"

wget -O jdk.tar.gz https://github.com/AdoptOpenJDK/openjdk8-binaries/releases/download/jdk8u222-b10/OpenJDK8U-jre_x64_linux_hotspot_8u222b10.tar.gz
tar -xzf jdk.tar.gz
rm jdk.tar.gz

./gradlew assemble
mv build/libs/bosh-plugin-sample.jar .

rm -rf "$ROOT/example/sample-plugin-java.zip"
zip -r "$ROOT/example/sample-plugin-java.zip" \
    "plugin.yml" \
    "bosh-plugin-sample.jar" \
    "java" \
    "jdk8u222-b10-jre"
rm -rf "jdk8u222-b10-jre"
rm bosh-plugin-sample.jar
popd


