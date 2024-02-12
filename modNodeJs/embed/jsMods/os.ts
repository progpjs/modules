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

// https://nodejs.org/api/os.html

interface ModOS {
    homeDir(): string
    hostName(): string
    tempDir(): string
}

interface ModProcess {
    arch(): string
    platform(): string
}

const modProcess = progpGetModule<ModProcess>("nodejsModProcess")!;
const modOS = progpGetModule<ModOS>("nodejsModOS")!;

export const arch = modProcess.arch;
export const machine = modProcess.arch;
export const platform = modProcess.platform;

export const homedir = modOS.homeDir;
export const hostname = modOS.hostName;
export const tmpdir = modOS.tempDir;

export const oseol = "\r";

export default {
    arch: arch,
    machine: machine,
    platform: platform,

    homedir: homedir,
    hostname: hostname,
    tmpdir: tmpdir,

    oseol: oseol
}