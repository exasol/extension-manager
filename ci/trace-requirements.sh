#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

oft_version="3.5.0"
tmp_dir="/tmp/oft/"
tmp_file="$tmp_dir/openfasttrace-$oft_version.jar"

base_dir="$( cd "$(dirname "$0")/.." >/dev/null 2>&1 ; pwd -P )"
readonly base_dir

if [[ ! -f "$tmp_file" ]]; then
    mkdir -p "$tmp_dir"
    url="https://repo1.maven.org/maven2/org/itsallcode/openfasttrace/openfasttrace/$oft_version/openfasttrace-$oft_version.jar"
    echo "Downloading $url to $tmp_file"
    curl --output "$tmp_file" "$url"
fi

# Trace all
#java -jar "$tmp_file" trace "$base_dir/doc"

# Trace only feat,req,dsn
java -jar "$tmp_file" trace --wanted-artifact-types feat,req,dsn "$base_dir/doc"