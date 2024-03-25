

export enum Stage { Test = "test", Prod = "prod" }

export function getAvailableStages(): string[] {
    return Object.keys(Stage).map(value => value.toLowerCase())
}

export function getExtensionGitHubRepos(): string[] {
    // "oracle-virtual-schema" extension is not yet released, see https://github.com/exasol/oracle-virtual-schema/issues/45
    return ["s3-document-files-virtual-schema", "row-level-security-lua", "cloud-storage-extension",
        "kafka-connector-extension", "kinesis-connector-extension"]
}

export interface CommandLineArgs {
    stage: Stage
    dryRun: boolean
}
