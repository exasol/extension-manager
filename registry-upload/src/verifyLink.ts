import { URL } from "node:url";


export async function verifyLink(url: URL) {
    try {
        await httpRequest(url)
    } catch (error) {
        throw new Error(`URL ${url.toString()} is invalid`, { cause: error })
    }
}

/**
 * Send an HTTP request and ensure it succeeds.
 * @param url the URL to request
 * @returns nothing
 * @throws Error if status is not OK or if the request failed
 */
async function httpRequest(url: URL): Promise<void> {
    const response = await fetch(url, { method: "HEAD" })
    if (!response.ok) {
        throw new Error(`Request HEAD ${url.toString()} failed with status ${response.status}`)
    }
}
