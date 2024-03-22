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

interface GitHubAsset {
    url: string;
    browser_download_url: string;
    name: string;
    content_type: string;
    size: number;
}

async function fetchLatestExtensions(gitHubRepos: string[]): Promise<Extension[]> {
    const octokit = await createGitHubClient()
    async function fetchLatestExtension(gitHubRepo: string): Promise<Extension> {
        const latestRelease = await octokit.rest.repos.getLatestRelease({ owner: "exasol", repo: gitHubRepo })
        const version = latestRelease.data.tag_name
        const extensionAssetName = getExtensionAssetName(latestRelease.data.assets, gitHubRepo, version);
        logAdapterJar(latestRelease.data.assets, gitHubRepo, version)
        return { id: gitHubRepo, url: buildExtensionUrl(gitHubRepo, version, extensionAssetName) }
    }
    return Promise.all(gitHubRepos.map(fetchLatestExtension))
}

function getExtensionAssetName(assets: GitHubAsset[], gitHubRepo: string, version: string): string {
    const extensionAssetsNames = assets.map(asset => asset.name).filter(name => name.endsWith(".js"));
    if (extensionAssetsNames.length !== 1) {
        throw new Error(`Expected exactly one .js extension in release ${getGitHubReleaseUrl(gitHubRepo, version)}, but got ${extensionAssetsNames.length}`);
    }
    const extensionAssetName = extensionAssetsNames[0];
    return extensionAssetName;
}

function getGitHubReleaseUrl(gitHubRepo: string, version: string) {
    return `https://github.com/exasol/${gitHubRepo}/releases/tag/${version}`;
}

function logAdapterJar(assets: GitHubAsset[], gitHubRepo: string, version: string) {
    const adapterAssets = assets.filter(asset => isAdapterAsset(asset.name));
    if (adapterAssets.length !== 1) {
        throw new Error(`Expected exactly one .jar adapter in release ${getGitHubReleaseUrl(gitHubRepo, version)}, but got ${adapterAssets.length}`);
    }
    console.log(`Adapter for ${gitHubRepo} ${version}: ${adapterAssets[0].browser_download_url}`);
}

function isAdapterAsset(assetName: string): boolean {
    if (assetName.endsWith(".lua")) {
        return true
    }
    return assetName.endsWith(".jar") && !assetName.includes("javadoc")
}

function buildExtensionUrl(gitHubRepo: string, version: string, extensionAssetName: string): string {
    // We use the internal URL to avoid inconsistent statistics in the public CDN
    return `https://extensions-internal.exasol.com/com.exasol/${gitHubRepo}/${version}/${extensionAssetName}`;
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
