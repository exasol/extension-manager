import { Stage } from "./common"

function printUsage() {
    const stages = Object.keys(Stage)
        .filter(v => typeof (v) === "string")
        .filter(v => v.length > 1)
        .map(v => v.toUpperCase())
        .join("|")
    console.error(`Usage: ts-node upload-registry-content.ts ${stages}`)
}

function getStage(): Stage | undefined {
    if (process.argv.length !== 3) {
        return undefined
    }
    const stageParam = process.argv[2].toLowerCase()
    switch (stageParam) {
        case "test": return Stage.Test
        case "prod": return Stage.Prod
        default: return undefined
    }
}

const stage = getStage()
if (stage === undefined) {
    printUsage()
    process.exit(1)
}
console.log("hello", stage)