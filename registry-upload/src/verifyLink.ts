import { FollowResponse, https } from "follow-redirects";
import { IncomingMessage } from "http";
import { RequestOptions } from "https";
import { URL } from "url";


export async function verifyLink(url: URL) {
    try {
        await httpRequest(url, { method: "HEAD" })
    } catch (error) {
        throw new Error(`URL ${url.toString()} is invalid`, { cause: error })
    }
}

/**
 * Send an HTTP request and return the response data.
 * <p>
 * This uses the `follow-redirects` library to follow redirects -- unlike Node's default `request()` function.
 * @param url the URL to request
 * @param options request options
 * @returns the response data
 * @throws Error if status is not OK or if the request failed
 */
function httpRequest(url: URL, options: RequestOptions): Promise<string> {
    // Convert the `request()` function to a Promise.
    // Callback `resolve()` signals a successful result, `reject()` signals a failure.
    return new Promise(function (resolve, reject) {
        const req = https.request(url, options, function (res: IncomingMessage & FollowResponse) {
            // Verify status code is 2xx OK
            if (res.statusCode && (res.statusCode < 200 || res.statusCode >= 300)) {
                return reject(new Error(`Request ${options.method} ${url.toString()} failed with status ${res.statusCode}`));
            }
            // Collect response data
            const chunks: Uint8Array[] = [];
            res.on('data', function (chunk: Uint8Array) {
                chunks.push(chunk);
            });
            // Build complete response and resolve the promise
            res.on('end', function () {
                try {
                    resolve(Buffer.concat(chunks).toString("utf-8"));
                } catch (e) {
                    reject(e);
                }
            });
        });
        // Handle request errors
        req.on('error', function (err: Error) {
            reject(err);
        });
        req.end();
    });
}
