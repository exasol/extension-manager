import { exec } from "child_process";
import { writeFileSync } from "fs";
import { Octokit } from "octokit";
import { promisify } from "util";
import { Stage, getExtensionGitHubRepos } from "./common";

const TESTING_EXTENSION: Extension = { id: "testing-extension", url: "https://d3d6d68cbkri8h.cloudfront.net/testing-extension.js" }

generateAllRegistries().catch((reason) => console.log("Failure:", reason))


async function generateAllRegistries() {
    const extensions = await fetchLatestExtensions(getExtensionGitHubRepos())
    const testStageExtensions = extensions.concat(TESTING_EXTENSION)
    writeRegistry(Stage.Prod, { extensions })
    writeRegistry(Stage.Test, { extensions: testStageExtensions })
}

interface RegistryContent {
    extensions: Extension[]
}

interface Extension {
    id: string
    url: string
}

async function fetchLatestExtensions(gitHubRepos: string[]): Promise<Extension[]> {
    const octokit = await createGitHubClient()
    async function fetchLatestExtension(gitHubRepo: string): Promise<Extension> {
        const latestRelease = await octokit.rest.repos.getLatestRelease({ owner: "exasol", repo: gitHubRepo })
        const version = latestRelease.data.tag_name
        const extensionAssetsNames = latestRelease.data.assets.map(asset => asset.name).filter(name => name.endsWith(".js"))
        if (extensionAssetsNames.length !== 1) {
            const url = `https://github.com/exasol/${gitHubRepo}/releases/tag/${version}`
            throw new Error(`Expected exactly one .js extension in release ${url}, but got ${extensionAssetsNames.length}`)
        }
        return { id: gitHubRepo, url: `https://extensions-internal.exasol.com/com.exasol/${gitHubRepo}/${version}/${extensionAssetsNames[0]}` }
    }
    return Promise.all(gitHubRepos.map(fetchLatestExtension))
}

function writeRegistry(stage: Stage, content: RegistryContent) {
    writeFileSync(`content/${stage}-registry.json`, JSON.stringify(content, null, 2))
}

async function createGitHubClient(): Promise<Octokit> {
    const token = await getGitHubToken()
    return new Octokit({ auth: token });
}

async function getGitHubToken(): Promise<string> {
    const asyncExec = promisify(exec)
    const { stdout } = await asyncExec("gh auth token")
    return stdout
}
