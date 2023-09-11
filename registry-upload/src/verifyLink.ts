import { request } from "follow-redirects/https";
import { RequestOptions } from "https";
import { URL } from "url";


export async function verifyLink(url: URL) {
    try {
        const content = await httpRequest(url, { method: "HEAD" })
    } catch (error) {
        throw new Error(`URL ${url} is invalid`, { cause: error })
    }
}

function httpRequest(url: string | URL, options: RequestOptions): Promise<string> {
    return new Promise(function (resolve, reject) {
        const req = request(url, options, function (res) {
            // reject on bad status
            if (res?.statusCode && (res.statusCode < 200 || res.statusCode >= 300)) {
                return reject(new Error(`Request ${options.method} ${url} failed with status ${res.statusCode}`));
            }
            // cumulate data
            const chunks: Uint8Array[] = [];
            res.on('data', function (chunk) {
                chunks.push(chunk);
            });
            res.on('end', function () {
                try {
                    resolve(Buffer.concat(chunks).toString());
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
