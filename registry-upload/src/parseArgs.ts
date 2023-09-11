import { ParseArgsConfig, parseArgs } from "node:util"
import { Stage } from "./common"


export interface CommandLineArgs {
    stage: Stage
    dryRun: boolean
}

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
    const config: ParseArgsConfig = { args, options, allowPositionals: false, strict: true }
    const { values } = parseArgs(config)
    console.log(`Args:`, args)
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
            throw new Error(`Got unexpected value ${value} for stage, allowed values: test, prod`)
    }
}

