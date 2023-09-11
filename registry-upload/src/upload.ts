import { existsSync } from "fs";
import { readFile } from "fs/promises";
import { resolve } from "path";
import { invalidateCloudFrontCache, readStackConfiguration, uploadFileContent } from "./awsService";
import { CommandLineArgs, Stage } from "./common";
import { verifyLink } from "./verifyLink";

const EXTENSION_MANAGER_STACK_NAME = "ExtensionManagerRegistry";

export async function upload(args: CommandLineArgs) {
    await verifyExtensionUrls(args.stage)
    const config = await readStackConfiguration(EXTENSION_MANAGER_STACK_NAME)
    console.log(`Read configuration: bucket ${config.bucketName}, domain: ${config.domainName}`)
    if (args.dryRun) {
        console.log("Dry run, skipping upload")
        return
    }
    await uploadFiles(args.stage, config.bucketName)
    await invalidateCloudFrontCache(config.cloudFrontDistributionId)
    await verifyLink(new URL(`https://${config.domainName}/registry.json`))
}

async function uploadFiles(stage: Stage, bucketName: string) {
    const promises: Promise<void>[] = []
    promises.push(uploadFileContent(bucketName, "registry.json", getRegistryFile(stage)))
    if (stage === Stage.Test) {
        promises.push(uploadFileContent(bucketName, "testing-extension.js", getTestingExtensionPath()))
    }
    console.log(`Uploading ${promises.length} files to S3...`)
    return Promise.all(promises)
}

function getTestingExtensionPath() {
    const filePath = resolve("../extension-manager-integration-test-java/testing-extension/dist/testing-extension.js")
    if (!existsSync(filePath)) {
        throw new Error(`Testing extension does not exist at '${filePath}'. Build it with "npm run build"`)
    }
    return filePath;
}

function getRegistryFile(stage: Stage) {
    const filePath = resolve(`content/${stage}-registry.json`)
    if (!existsSync(filePath)) {
        throw new Error(`Registry file content does not exist at '${filePath}'`)
    }
    return filePath;
}

async function verifyExtensionUrls(stage: Stage): Promise<void> {
    const path = getRegistryFile(stage)
    const content = await readFile(path, { encoding: "utf-8" })
    const registry: any = JSON.parse(content)
    console.log(`Verify links for ${registry.extensions.length} extensions...`)
    await Promise.all(registry.extensions.map((ext: any) => verifyExtensionEntry(ext)))
    console.log(`All ${registry.extensions.length} links are valid.`)
}

async function verifyExtensionEntry(extension: any): Promise<void> {
    const id: string = extension.id
    const url: string = extension.url
    try {
        const content = await verifyLink(new URL(url))
    } catch (error) {
        throw new Error(`URL for extension ${id} is invalid: ${error}`, { cause: error })
    }
}


