#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

stage=${1:-}

function verify_arguments() {
    if [ -z "$stage" ]; then
        echo "Stage to deploy not specified."
        echo "Usage: $0 test|prod"
        exit 1
    fi

    registry_file="content/$stage-registry.json"
    if [ ! -f "$registry_file" ]; then
        echo "Registry file '$registry_file' does not exist for specified stage '$stage'"
        exit 1
    fi
}

function get_cloudformation_info() {
    stack_output=$(aws cloudformation describe-stacks --stack-name ExtensionManagerRegistry --output=text --query "Stacks[0].Outputs[].{key:ExportName,value:OutputValue}")
    bucket_name="$(echo "$stack_output" | grep StaticContentBucketName | awk -F '\t' '{print $2}')"
    cloudfront_distribution_id="$(echo "$stack_output" | grep CloudFrontDistributionId | awk -F '\t' '{print $2}')"
    domain_name="$(echo "$stack_output" | grep CloudFrontDistributionDomainName | awk -F '\t' '{print $2}')"
    registry_url="https://$domain_name/registry.json"
    echo "S3 Bucket: $bucket_name, CloudFront distribution $cloudfront_distribution_id, URL: $registry_url"
}

function upload_registry() {
    local s3_url="s3://$bucket_name/registry.json"
    echo "Uploading $registry_file to $s3_url"
    aws s3 cp "$registry_file" "$s3_url"
}

function upload_testing_extension() {
    if [ "$stage" != "test" ]; then
        echo "Skip uploading testing extension"
        return
    fi
    local file_name="testing-extension.js"
    local testing_extension="../extension-manager-integration-test-java/testing-extension/dist/$file_name"
    if [ ! -f "$testing_extension" ]; then
        echo "Testing extension '$testing_extension' does not exist. Build it first with 'npm run build'"
        exit 1
    fi
    local s3_url="s3://$bucket_name/$file_name"
    echo "Uploading $testing_extension to $s3_url"
    aws s3 cp "$testing_extension" "$s3_url"
}

function invalidate_cloudfront_cache() {
    echo "Invalidating cloudfront cache for distribution $cloudfront_distribution_id"
    aws cloudfront create-invalidation --distribution-id "$cloudfront_distribution_id" --paths '/*'
}

verify_arguments

echo "Using AWS_PROFILE=$AWS_PROFILE to deploy $registry_file"

get_cloudformation_info

upload_registry

upload_testing_extension

invalidate_cloudfront_cache