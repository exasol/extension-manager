#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

oft_version="3.6.0"
tmp_dir="/tmp/oft/"
jar_file="$tmp_dir/openfasttrace-$oft_version.jar"

base_dir="$( cd "$(dirname "$0")/.." >/dev/null 2>&1 ; pwd -P )"
readonly base_dir

if [[ ! -f "$jar_file" ]]; then
    mkdir -p "$tmp_dir"
    url="https://repo1.maven.org/maven2/org/itsallcode/openfasttrace/openfasttrace/$oft_version/openfasttrace-$oft_version.jar"
    echo "Downloading $url to $jar_file"
    curl --output "$jar_file" "$url"
fi

# Trace all
java -jar "$jar_file" trace \
    "$base_dir/doc" \
    "$base_dir/pkg" \
    "$base_dir/extension-manager-integration-test-java" \
    --wanted-artifact-types impl,itest,utest,dsn
