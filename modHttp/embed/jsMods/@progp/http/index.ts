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
    requestRemainingSegments(resId: SharedResource): string[]|null;
    requestCookie(resId: SharedResource, name: string): HttpCookie|null;
    requestHeaders(resId: SharedResource): any;
    requestCookies(resId: SharedResource): {[key:string]:HttpCookie};

    requestReadFormFile(resId: SharedResource, fieldName: string, fileId: number, callback: Function): any;
    requestSaveFormFile(resId: SharedResource, fieldName: string, fileId: number, saveFilePath: string, callback: Function): void;
}

interface CookieOptions {
    isHttpOnly?:   boolean
    isSecure?:     boolean
    sameSiteType?: number
    domaine?:      string
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
    private _requestRemainingSegments: string[]|null|undefined;
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

    requestRemainingSegments(): any {
        if (this._requestRemainingSegments===undefined) {
            return this._requestRemainingSegments = modHttp.requestRemainingSegments(this.resId);
        }

        return this._requestRemainingSegments;
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

export interface HttpServerConfig {
}

export class HttpServer {
    private readonly serverPort: number;
    private isStarted: boolean = false;

    constructor(serverPort: number) {
        this.serverPort = serverPort;
    }

    configure(config: HttpServerConfig) {
        // Return false if the server is already started.
        if (!modHttp.configureServer(this.serverPort, config)) {
            this.isStarted = true;
        }
    }

    start() {
        if (this.isStarted) return;
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
}

export type HttpRequestHandler = (res: HttpRequest) => Promise<void>;

export function asHttpRequest(f: (req:HttpRequest)=>void) {
    return (resId:SharedResource)=> f(new HttpRequest(resId, gSecureCaller))
}

let gSecureCaller = {v: {}};
