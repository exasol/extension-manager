import { CommandLineArgs, getAvailableStages } from "./common"
import { parseArguments } from "./parseArgs"
import { upload } from "./upload"

function printUsage() {
    const stages = getAvailableStages().join("|")
    console.error(`Usage: AWS_PROFILE=$profile npm run upload -- --stage=${stages} [--no-dry-run]`)
}

function getArgs(): CommandLineArgs {
    const allArgs = process.argv
    const relevantArgs = allArgs.splice(2, allArgs.length - 1)
    try {
        return parseArguments(relevantArgs)
    } catch (error) {
        console.error(`Failed to parse arguments [${relevantArgs.join(' ')}]:`, error)
        printUsage()
        process.exit(1)
    }
}

const args = getArgs()
upload(args).catch((reason) => console.log("Failure:", reason))
