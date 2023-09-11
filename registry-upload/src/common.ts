

export enum Stage { Test = "test", Prod = "prod" }

export function getAvailableStages(): string[] {
    return Object.keys(Stage).map(value => value.toLowerCase())
}

export interface CommandLineArgs {
    stage: Stage
    dryRun: boolean
}

