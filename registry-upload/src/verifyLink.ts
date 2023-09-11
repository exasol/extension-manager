import { request } from "follow-redirects/https";
import { RequestOptions } from "https";
import { URL } from "url";


export async function verifyLink(url: URL) {
    try {
        await httpRequest(url, { method: "HEAD" })
    } catch (error) {
        throw new Error(`URL ${url.toString()} is invalid`, { cause: error })
    }
}

function httpRequest(url: URL, options: RequestOptions): Promise<string> {
    return new Promise(function (resolve, reject) {
        const req = request(url, options, function (res) {
            if (res.statusCode && (res.statusCode < 200 || res.statusCode >= 300)) {
                return reject(new Error(`Request ${options.method} ${url.toString()} failed with status ${res.statusCode}`));
            }
            const chunks: Uint8Array[] = [];
            res.on('data', function (chunk: Uint8Array) {
                chunks.push(chunk);
            });
            res.on('end', function () {
                try {
                    resolve(Buffer.concat(chunks).toString("utf-8"));
                } catch (e) {
                    reject(e);
                }
            });
        });
        req.on('error', function (err) {
            reject(err);
        });
        req.end();
    });
}
