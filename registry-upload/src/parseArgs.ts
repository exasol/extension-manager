import { ParseArgsConfig, parseArgs } from "node:util"
import { CommandLineArgs, Stage, getAvailableStages } from "./common.js"

export function parseArguments(args: string[]): CommandLineArgs {
    const options: any = {
        stage: {
            type: "string"
        },
        "no-dry-run": {
            type: "boolean",
            default: false
        }
    }
    // eslint-disable-next-line @typescript-eslint/no-unsafe-assignment
    const config: ParseArgsConfig = { args, options, allowPositionals: false, strict: true }
    const { values } = parseArgs(config)
    return {
        stage: parseStage(values.stage),
        dryRun: values["no-dry-run"] === false
    }
}

function parseStage(value: any): Stage {
    if (typeof value !== "string") {
        throw new Error(`Type of stage ${value} is ${typeof value}, expected string`)
    }
    const stageParam = value.toLowerCase()
    switch (stageParam) {
        case "test": return Stage.Test
        case "prod": return Stage.Prod
        default:
            throw new Error(`Got unexpected value ${value} for stage, allowed values: ${getAvailableStages().join(", ")}`)
    }
}
