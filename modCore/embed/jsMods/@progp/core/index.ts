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

declare global {
    /**
     * Get a natif function group (function in C++ or Go).
     * Return undefined if the security level is inferior to the security level of this group.
     */
    function progpGetModule<T>(modName: string): T | undefined; // Declared in C++, group "progpCore".

    /**
     * Print his argument to the console.
     * Is like "console.log".
     */
    function progpPrint(...params: any): void;  // Declared in C++, group "progpCore".

    /**
     * Allows to avoid that a function is removed by the javascript optimizer.
     */
    function progpDontRemove(e: any): void; // Declared in this js script.

    function progpCallAfterMs(timeInMs: number, callback: Function): void;

    function progpStringToBuffer(text: string): ArrayBuffer;
    function progpBufferToString(b: ArrayBuffer): string;

    /**
     * progpDispose will dispose the SharedResource.
     * It will call the dispose function binded to this resource
     * and dispose the resource wrapper.
     */
    function progpDispose(res: SharedResource): void;

    /**
     * progpAutoDispose will automatically dispose all the resources created
     * by the f function once this function is fully terminated, which include
     * all the asynchrone function call from this function.
     *
     * For exemple if you open a file, this file resource will be disposed
     * once this function ends, without having to manually call progpDispose(myFile).
     */
    function progpAutoDispose(f: Function): void;
}

//region Web Standard libraries

//region Timers

let g_nextTimerId = 1;
let g_timers: {[key: number]: boolean} = [];

// @ts-ignore
globalThis.setTimeout = function (callbackFct: Function, timeInMs: number, ...params: any): number {
    if (!callbackFct) return -1;
    const timerId = g_nextTimerId++;
    g_timers[timerId] = true;

    progpCallAfterMs(timeInMs, () => {
        let timerState = g_timers[timerId];

        if (timerState) {
            delete(g_timers[timerId]);
            callbackFct.call(globalThis, params);
        }
    });

    return timerId;
};

function setIntervalAux(callbackFct: Function, timeInMs: number, params: any[], timerId: number) {
    progpCallAfterMs(timeInMs, () => {
        let timerState = g_timers[timerId];

        if (timerState) {
            callbackFct.call(globalThis, params);
            setIntervalAux(callbackFct, timeInMs, params, timerId);
        }
    });
}

// @ts-ignore
globalThis.setInterval = function (callbackFct: Function, timeInMs: number, ...params: any): number {
    if (!callbackFct) return -1;
    const timerId = g_nextTimerId++;
    g_timers[timerId] = true;

    setIntervalAux(callbackFct, timeInMs, params, timerId);
    return timerId;
}

// @ts-ignore
globalThis.clearTimeout = function(timerId: number) {
    delete(g_timers[timerId]);
};

globalThis.clearInterval = globalThis.clearTimeout;

//endregion

//region Console

const bckConsoleLog = globalThis.console.log;
const bckConsoleWarn = globalThis.console.warn;
const bckConsoleError = globalThis.console.error;
const bckConsoleDebug = globalThis.console.debug;
const bckConsoleInfo = globalThis.console.info;

globalThis.console.log = function(...data: any[]) {
    progpPrint(...data);
    if (bckConsoleLog) bckConsoleLog(...data);
}

globalThis.console.warn = function(...data: any[]) {
    progpPrint("[WARN] ", ...data);
    if (bckConsoleWarn) bckConsoleWarn(...data);
}

globalThis.console.error = function(...data: any[]) {
    progpPrint("[ERROR] ", ...data);
    if (bckConsoleError) bckConsoleError(...data);
}

globalThis.console.debug = function(...data: any[]) {
    progpPrint("[DEBUG] ", ...data);
    if (bckConsoleDebug) bckConsoleDebug(...data);
}

globalThis.console.info = function(...data: any[]) {
    progpPrint("[INFO] ", ...data);
    if (bckConsoleInfo) bckConsoleInfo(...data);
}

//endregion

//endregion

//region NodeJS buffer

const nodeJsBuffer = require("node:buffer");
globalThis.Buffer = nodeJsBuffer;

//endregion

/**
 * Allows to known when a SharedResource is required.
 * Here it's mainly a placeholder since the underlying type is a number.
 */
export interface SharedResource {}
