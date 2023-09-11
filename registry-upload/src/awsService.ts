import { CloudFormationClient, DescribeStacksCommand, Output, Stack } from "@aws-sdk/client-cloudformation";
import { CloudFrontClient, CreateInvalidationCommand } from "@aws-sdk/client-cloudfront";
import { S3 } from "@aws-sdk/client-s3";
import { createReadStream } from "fs";

const AWS_REGION = "eu-central-1";

export interface StackConfig {
    bucketName: string
    cloudFrontDistributionId: string
    domainName: string
}

const cloudFormationClient = new CloudFormationClient({ region: AWS_REGION })
const s3Client = new S3({ region: AWS_REGION })
const cloudFrontClient = new CloudFrontClient({ region: AWS_REGION })

async function describeStack(stackName: string): Promise<Stack> {
    const stacks = await cloudFormationClient.send(new DescribeStacksCommand({ StackName: stackName }))
    if (stacks.Stacks === undefined) {
        throw new Error(`Failed to describe stack ${stackName}`)
    }
    if (stacks.Stacks.length !== 1) {
        throw new Error(`Expected exactly one stack with name ${stackName} but got ${stacks.Stacks.length}`)
    }
    return stacks.Stacks[0]
}

async function getStackOutputs(stackName: string): Promise<Output[]> {
    const stack = await describeStack(stackName)
    if (stack.Outputs === undefined) {
        throw new Error(`Stack '${stackName}' does not have outputs`)
    }
    return stack.Outputs
}

export async function readStackConfiguration(cloudFormationStackName: string): Promise<StackConfig> {
    const outputs = await getStackOutputs(cloudFormationStackName)

    function getOutputValue(name: string): string {
        const output = outputs.filter(o => o.ExportName === name)
        if (output.length !== 1) {
            throw new Error(`Expected exactly one output named ${name} but found ${output.length} for stack ${cloudFormationStackName}`)
        }
        if (output[0].OutputValue === undefined) {
            throw new Error(`Output ${name} has no value for stack ${cloudFormationStackName}`)
        }
        return output[0].OutputValue
    }

    return {
        bucketName: getOutputValue("StaticContentBucketName"),
        cloudFrontDistributionId: getOutputValue("CloudFrontDistributionId"),
        domainName: getOutputValue("CloudFrontDistributionDomainName")
    }
}

export async function uploadFileContent(bucket: string, key: string, localFilePath: string): Promise<void> {
    const stream = createReadStream(localFilePath)
    await s3Client.putObject({ Bucket: bucket, Key: key, Body: stream })
}

export async function invalidateCloudFrontCache(distributionId: string): Promise<void> {
    console.log(`Invalidating CloudFront cache for distribution ${distributionId}...`)
    const callerReference = Date.now().toString()
    const invalidationBatch = { CallerReference: callerReference, Paths: { Quantity: 1, Items: ["/*"] } };
    await cloudFrontClient.send(new CreateInvalidationCommand({ DistributionId: distributionId, InvalidationBatch: invalidationBatch }))
}