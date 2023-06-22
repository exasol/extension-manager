#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

oft_version="3.7.0"

base_dir="$( cd "$(dirname "$0")/.." >/dev/null 2>&1 ; pwd -P )"
readonly base_dir
readonly oft_jar="$HOME/.m2/repository/org/itsallcode/openfasttrace/openfasttrace/$oft_version/openfasttrace-$oft_version.jar"

if [ ! -f "$oft_jar" ]; then
    echo "Downloading OpenFastTrace $oft_version"
    mvn --batch-mode org.apache.maven.plugins:maven-dependency-plugin:3.3.0:get -Dartifact=org.itsallcode.openfasttrace:openfasttrace:$oft_version
fi

# Trace all
java -jar "$oft_jar" trace \
    "$base_dir/doc" \
    "$base_dir/pkg" \
    "$base_dir/extension-manager-integration-test-java"
