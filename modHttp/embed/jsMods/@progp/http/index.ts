/*
 * (C) Copyright 2024 Johan Michel PIQUET, France (https://johanpiquet.fr/).
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// noinspection JSUnusedGlobalSymbols

import {SharedResource} from "@progp/core"

//region Go plugin

interface ModHttpServer {
    startServer(serverPort: number): void;
    configureServer(serverPort: number, config: any): boolean;

    getHost(serverPort: number, hostName: string): SharedResource
    VERB_withFunction(hostRes: SharedResource, verb: string, requestPath: string, handler: Function): void
    
    returnString(resId: SharedResource, httpCode: number, contentType: string, value: string): void;
    responseSetHeader(resId: SharedResource, key: string, value: string): void;
    responseSetCookie(resId: SharedResource, key: string, value: string, options: CookieOptions): void;

    requestURI(resId: SharedResource): string;
    requestPath(resId: SharedResource): string;
    requestIP(resId: SharedResource): string;
    requestMethod(resId: SharedResource): string;
    requestHost(resId: SharedResource): string;
    requestQueryArgs(resId: SharedResource): any;
    requestPostArgs(resId: SharedResource): any;
    requestWildcards(resId: SharedResource): string[]|null;
    requestCookie(resId: SharedResource, name: string): HttpCookie|null;
    requestHeaders(resId: SharedResource): any;
    requestCookies(resId: SharedResource): {[key:string]:HttpCookie};

    requestReadFormFile(resId: SharedResource, fieldName: string, fileId: number, callback: Function): any;
    requestSaveFormFile(resId: SharedResource, fieldName: string, fileId: number, saveFilePath: string, callback: Function): void;

    sendFile(resId: SharedResource, filePath: string): void
    sendFileAsIs(resId: SharedResource, filePath: string, mimeType: string, contentEncoding: string): void

    gzipCompressFile(sourceFile: string, destFile: string, compressionLevel: number, callback: Function): void;
    brotliCompressFile(sourceFile: string, destFile: string, compressionLevel: number, callback: Function): void;

    fetch(url: string, options: FetchOptions, callback: Function): void;

    proxyTo(resId: SharedResource, fromPath: string, targetHost: string, options: ProxyTypeOptions): void

    fileServer_Create(resId: SharedResource, fromPath: string, dirPath: string, options: ServeFileOptions): SharedResource
    fileServer_RemoveAll(resId: SharedResource): void
    fileServer_RemoveUri(resId: SharedResource, uri: string, data: string): void
    fileServer_VisitCache(resId: SharedResource, callback: Function): void
    fileServer_OnFileNotFound(resId: SharedResource, callback: Function): void
}

interface CookieOptions {
    isHttpOnly?:   boolean
    isSecure?:     boolean
    sameSiteType?: number
    domain?:       string
    expireTime?:   number
    maxAge?:       number
}

const modHttp = progpGetModule<ModHttpServer>("progpjsModHttp")!;

//endregion

export interface HttpCookie {
    key: string;
    domain: string;
    value: string;
    maxAge: number;
    expireTime: number;
    sameSiteType: number;
    isSecure: boolean;
    isHTTPOnly: boolean;
}

export interface FormFileInfo {
    /* Offset of this file inside the list of files returned for this form field */
    id: number

    header: any
    size: number
}

export class HttpRequest {
    private readonly resId: SharedResource;
    private _requestURI: string|undefined;
    private _requestPath: string|undefined;
    private _requestIP: string|undefined;
    private _requestMethod: string|undefined;
    private _requestHost: string|undefined;
    private _requestQueryArgs: any|undefined;
    private _requestPostArgs: any|undefined;
    private _requestWildcards: string[]|null|undefined;
    private _requestHeaders: any|undefined;
    private _requestCookies: {[key:string]:HttpCookie}|undefined;
    private _contentType: string = "text/html";

    constructor(resId: SharedResource, caller: any) {
        if (caller!==gSecureCaller) throw Error("Forbiden call");
        this.resId = resId;
    }

    returnString(httpCode: number, value: string) {
        modHttp.returnString(this.resId, httpCode, this._contentType, value);
    }

    returnHtml(httpCode: number, value: string) {
        modHttp.returnString(this.resId, httpCode, this._contentType, value);
    }

    sendFile(filePath: string) {
        modHttp.sendFile(this.resId, filePath);
    }

    sendFileAsIs(filePath: string, mimeType?: string, contentEncoding?: string) {
        if (mimeType===undefined) mimeType = "";
        if (contentEncoding===undefined) contentEncoding = "";
        modHttp.sendFileAsIs(this.resId, filePath, mimeType, contentEncoding);
    }

    setHeader(key: string, value: string) {
        modHttp.responseSetHeader(this.resId, key, value);
    }

    setContentType(value: string) {
        this._contentType = value;
    }

    deleteCookie(key: string) {
        this.setCookie(key, "",{maxAge: 0})
    }

    setCookie(key: string, value: string, options?: CookieOptions) {
        modHttp.responseSetCookie(this.resId, key, value, options||{});
    }

    requestURI(): string {
        if (this._requestURI===undefined) {
            return this._requestURI = modHttp.requestURI(this.resId);
        }

        return this._requestURI;
    }

    requestPath(): string {
        if (this._requestPath===undefined) {
            return this._requestPath = modHttp.requestPath(this.resId);
        }

        return this._requestPath;
    }

    requestIP(): string {
        if (this._requestIP===undefined) {
            return this._requestIP = modHttp.requestIP(this.resId);
        }

        return this._requestIP;
    }

    requestMethod(): string {
        if (this._requestMethod===undefined) {
            return this._requestMethod = modHttp.requestMethod(this.resId);
        }

        return this._requestMethod;
    }

    requestHost(): string {
        if (this._requestHost===undefined) {
            return this._requestHost = modHttp.requestHost(this.resId);
        }

        return this._requestHost;
    }

    /**
     * PHP like style, allows making thing easyier.
     */
    _$GET(key: string) {
        if (this._requestQueryArgs===undefined) {
            return this._requestQueryArgs = modHttp.requestQueryArgs(this.resId);
        }

        if (!this._requestQueryArgs) return undefined;
        return this._requestQueryArgs[key];
    }

    requestQueryArgs(): {[key:string]:string} {
        if (this._requestQueryArgs===undefined) {
            return this._requestQueryArgs = modHttp.requestQueryArgs(this.resId);
        }

        return this._requestQueryArgs;
    }

    readFormFile(key: string, fileOffset: number): Promise<ArrayBuffer|null> {
        return new Promise<ArrayBuffer|null>((resolve, reject) => {
            modHttp.requestReadFormFile(this.resId, key, fileOffset, (err: string, buffer: ArrayBuffer) => {
                if (err) reject(err);
                else resolve(buffer);
            });
        })
    }

    saveFormFile(key: string, fileOffset: number, saveFilePath: string): Promise<boolean> {
        return new Promise<boolean>((resolve, reject) => {
            modHttp.requestSaveFormFile(this.resId, key, fileOffset, saveFilePath, (err: string, ok: boolean) => {
                if (err) reject(err);
                else resolve(ok);
            });
        });
    }

    /**
     * PHP like style, allows making thing easyier.
     */
    _$POST(key: string) {
        if (this._requestPostArgs===undefined) {
            return this._requestPostArgs = modHttp.requestPostArgs(this.resId);
        }

        if (!this._requestPostArgs) return undefined;
        return this._requestPostArgs[key];
    }

    requestPostArgs(): {[key:string]: string | FormFileInfo[]} {
        if (this._requestPostArgs===undefined) {
            return this._requestPostArgs = modHttp.requestPostArgs(this.resId);
        }

        return this._requestPostArgs;
    }

    requestWildcards(): any {
        if (this._requestWildcards===undefined) {
            return this._requestWildcards = modHttp.requestWildcards(this.resId);
        }

        return this._requestWildcards;
    }

    requestCookies(): {[key:string]:HttpCookie} {
        if (this._requestCookies===undefined) {
            let res = modHttp.requestCookies(this.resId);
            if (!res) res = {};
            return this._requestCookies = res;
        }

        return this._requestCookies;
    }

    requestCookie(name: string): HttpCookie|null {
        return this.requestCookies()[name];
    }

    requestHeaders(): any {
        if (this._requestHeaders===undefined) {
            return this._requestHeaders = modHttp.requestHeaders(this.resId);
        }

        return this._requestHeaders;
    }
}

export interface HttCertificate {
    hostName: string

    certFilePath?: string
    keyFilePath?: string

    useLetsEncrypt?: boolean,
    cacheDir?: string
}

export interface HttpServerConfig {
    /**
     * Allows to hide server error and don't write message in the console.
     * It's important to hide error, without what attackers can slow down the server by a lot with fake calls.
     */
    hideErrors?: boolean
    enableHttps?: boolean
    certificates: HttCertificate[]
}

export interface FetchOptions {
    /**
     * Indicate the http method to use.
     * Default is GET.
     */
    method?: string

    /**
     * If set, save the body inside a file.
     */
    streamBodyToFile?: string

    /**
     * If true, returns the response headers.
     */
    returnHeaders?: boolean

    /**
     * If true, returns the response cookies.
     */
    returnCookies?: boolean

    /**
     * If set, set the list of the headers to send.
     */
    sendHeaders?: {[key:string]:string}


    /**
     * If set, set the list of the cookies to send.
     */
    setCookies?: {[key:string]:string}

    /**
     * Allows to return body event if response code isn't 200 Ok.
     */
    forceReturningBody?: boolean

    /**
     * Allows to avoid requesting the body.
     * Is useful if testing target existence or when we want to only get his headers.
     */
    skipBody?: boolean

    /**
     * Set the content type used when sending a body with the request.
     * Isn't set when no request are set.
     */
    contentType?: string

    /**
     * UserAgent set the user agent used when sending a body with the request.
     * Isn't set when no request are set.
     */
    userAgent?: string
}

export interface ProxyTypeOptions {
    /**
     * Allows to avoid including url from this path.
     */
    excludeSubPaths?: boolean
}

export interface ServeFileOptions {

}

export class HttpServer {
    private readonly serverPort: number;
    private isStarted: boolean = false;
    private config: HttpServerConfig|undefined;

    constructor(serverPort: number) {
        this.serverPort = serverPort;
    }

    configure(config: HttpServerConfig) {
        this.config = config
    }

    /**
     * Register a https certificate.
     *
     * To create your dev certificate:
     * 1- Install mkcert from https://github.com/FiloSottile/mkcert
     * 2- mkcert -install		--> to do one time, create the root CA certificate, which is required to create your own test certificate.
     * 3- mkcert myhostname     --> create a valid certificate (here you can replace myhostname by localhost).
     *
     * Automatic dev certificate isn't supported anymore since it introduce difficult behaviors.
     */
    addHttpsCertificate(hostName: string, certFilePath: string, keyFilePath: string) {
        if (!this.config) {
            this.config = {
                enableHttps: true,
                certificates: []
            }
        }

        this.config!.certificates.push({
            hostName: hostName,
            certFilePath: certFilePath,
            keyFilePath: keyFilePath
        });
    }

    addLetEncryptCertificate(hostName: string, cacheDir: string) {
        if (!this.config) {
            this.config = {
                enableHttps: true,
                certificates: []
            }
        }

        this.config!.certificates.push({
            hostName: hostName,
            useLetsEncrypt: true,
            cacheDir: cacheDir
        });
    }

    start() {
        if (this.isStarted) return;

        if (this.config && !modHttp.configureServer(this.serverPort, this.config)) {
            // configureServer returns false if the server is already started.
            this.isStarted = true;
            return;
        }

        modHttp.startServer(this.serverPort);
        this.isStarted = true;
    }

    getHost(hostName: string): HttpHost {
        let hostResId = modHttp.getHost(this.serverPort, hostName);
        return new HttpHost(hostResId);
    }
}

export class HttpHost {
    private readonly hostResId: SharedResource;

    constructor(hostResId: SharedResource) {
        this.hostResId = hostResId;
    }

    verb(verb: string, requestPath: string, handler: HttpRequestHandler): void {
        modHttp.VERB_withFunction(this.hostResId, verb, requestPath, (_: string, resId: SharedResource) => {
            handler(new HttpRequest(resId, gSecureCaller));
        });
    }

    GET(requestPath: string, handler: HttpRequestHandler): void {
        this.verb("GET", requestPath, handler);
    }

    POST(requestPath: string, handler: HttpRequestHandler): void {
        this.verb("POST", requestPath, handler);
    }

    proxyTo(fromPath: string, targetHost: string, options?: ProxyTypeOptions) {
        if (!fromPath) fromPath = "/";
        if (!options) options = {};

        modHttp.proxyTo(this.hostResId, fromPath, targetHost, options)
    }

    serveFiles(fromPath: string, dirPath: string, options?: ServeFileOptions) {
        if (!fromPath) fromPath = "/";
        if (!options) options = {};

        let res = modHttp.fileServer_Create(this.hostResId, fromPath, dirPath, options)
        return new FileServer(res)
    }
}

export interface FsCacheEntry {
    hitCount: number
    filePath: string
    gzipFilePath: string
    uri: string
    data: string
    fileUpdateDate: number
    lastRequestedData: number
    contentType: string
    contentLength: number
    gzipContentLength: number
}

export interface FsOnFileNodeFound {
    filePath: string
    uri: string
    uriPath: string
    queryString: string
    data: string
    hostName: string
}

export class FileServer {
    private readonly resId: SharedResource

    constructor(resId: SharedResource) {
        this.resId = resId
    }

    removeAll() {
        modHttp.fileServer_RemoveAll(this.resId)
    }

    visitCache(onCacheEntry: ((v: FsCacheEntry) => void)) {
        modHttp.fileServer_VisitCache(this.resId, (_: string, raw: string) => {
            onCacheEntry(JSON.parse(raw))
        })
    }

    onFileNotFound(f: ((v: FsOnFileNodeFound, oneDone: Function) => void)) {
        modHttp.fileServer_OnFileNotFound(this.resId, (resId: SharedResource, json: string) => {
            f(<FsOnFileNodeFound>JSON.parse(json), () => progpReturnVoid(resId))
        })
    }

    /**
     * Remove the exact uri.
     * @param uri   The uri to remove.
     * @param data  Is used by hooks to finely select content to remove.
     */
    removeUri(uri: string, data?: string) {
        if (data===undefined) data = "";
        modHttp.fileServer_RemoveUri(this.resId, uri, data)
    }
}

export type HttpRequestHandler = (res: HttpRequest) => Promise<void>;

export function asHttpRequest(f: (req:HttpRequest)=>void) {
    return (resId:SharedResource)=> f(new HttpRequest(resId, gSecureCaller))
}

export const GZIP_COMPRESSION_LEVEL_NO_COMPRESSION = 0
export const GZIP_COMPRESSION_LEVEL_BEST_SPEED = 1
export const GZIP_COMPRESSION_LEVEL_BEST_COMPRESSION = 9
export const GZIP_COMPRESSION_LEVEL_DEFAULT_COMPRESSION = -1
export const GZIP_COMPRESSION_LEVEL_HUFFMAN_ONLY = -2

/**
 * Allows to compress a file with GZip in the same way as what http server do.
 * This allows to store a pre-compressed file in order to serve it as-is for full-speed.
 */
export async function gzipCompressFile(sourceFilePath: string, targetFilePath: string, compressionLevel?: number): Promise<void> {
    // Default is best compression level.
    if (!compressionLevel) compressionLevel = GZIP_COMPRESSION_LEVEL_BEST_COMPRESSION;

    return new Promise<void>(function (resolve, reject) {
        modHttp.gzipCompressFile(sourceFilePath, targetFilePath, compressionLevel!, (err: string)=>{
            if (err) reject(err);
            else resolve();
        });
    });
}

export const BROTLI_COMPRESSION_LEVEL_BEST_SPEED = 0
export const BROTLI_COMPRESSION_LEVEL_BEST_COMPRESSION = 11
export const BROTLI_COMPRESSION_LEVEL_DEFAULT_COMPRESSION = 6

/**
 * Allows to compress a file with Brotli in the same way as what http server do.
 * This allows to store a pre-compressed file in order to serve it as-is for full-speed.
 */
export async function brotliCompressFile(sourceFilePath: string, targetFilePath: string, compressionLevel?: number): Promise<void> {
    // Default is best compression level.
    if (!compressionLevel) compressionLevel = BROTLI_COMPRESSION_LEVEL_BEST_COMPRESSION;

    return new Promise<void>(function (resolve, reject) {
        modHttp.brotliCompressFile(sourceFilePath, targetFilePath, compressionLevel!, (err: string)=>{
            if (err) reject(err);
            else resolve();
        });
    });
}

export interface FetchResult {
    statusCode: number,
    body?: string
    headers?: {[key:string]: string}
    cookies?: {[key:string]: any}
}

export async function fetch(url: string, options?: FetchOptions): Promise<FetchResult> {
    if (!options) options = {};
    if (!options.method) options.method = "GET";

    return new Promise<FetchResult>(function (resolve, reject) {
        modHttp.fetch(url, options!, (err: string, res: string) => {
            if (err) reject(err);
            else resolve(JSON.parse(res));
        })
    })
}

let gSecureCaller = {v: {}};
