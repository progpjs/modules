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

// https://nodejs.org/api/fs.html

/*
    [e] Sync API
        [x] fs.existsSync(path)
        [x] fs.statSync(path[, options])
        [x] fs.accessSync(path[, mode])
        [x] fs.chmodSync(path, mode)
        [x] fs.chownSync(path, uid, gid)
        [e] fs.truncateSync(path[, len])

        [e] fs.copyFileSync(src, dest[, mode])
        [x] fs.linkSync(existingPath, newPath)                      --> need buffer
        [x] fs.symlinkSync(target, path[, type])
        [x] fs.unlinkSync(path)

        [x] fs.mkdirSync(path[, options])
        [x] fs.mkdtempSync(prefix[, options])
        [x] fs.renameSync(oldPath, newPath)
        [x] fs.rmdirSync(path[, options])
        [x] fs.rmSync(path[, options])

        [x] fs.readFileSync(path[, options])                        --> need buffer
        [ ] fs.appendFileSync(path, data[, options])                --> need buffer

        [ ] fs.readlinkSync(path[, options])                        --> need buffer
        [ ] fs.realpathSync(path[, options])                        --> need buffer

        [ ] fs.openSync(path[, flags[, mode]])                      --> need buffer
        [ ] fs.closeSync(fd)                                        --> need buffer
        [ ] fs.opendirSync(path[, options])                         --> need buffer
        [ ] fs.readSync(fd, buffer, offset, length[, position])     --> need buffer
        [ ] fs.writeFileSync(file, data[, options])                 --> need buffer
        [ ] fs.writeSync(fd, buffer, offset[, length[, position]])  --> need buffer
 */

interface ModFS {
    existsSync(path: string): boolean
    statSync(path: string, throwErrorIfMissing: boolean): any
    accessSync(path: string, mode: number): any
    chmodSync(path: string, mode: number): void
    chownSync(path: string, uid: number, gid: number): void
    truncateSync(path: string, length: number): void
    readFileUtf8Sync(path: string): string
    readFileBytesSync(path: string): ArrayBuffer
    copyFileSync(fromPath: string, toPath: string): void
    linkSync(existingPath: string, newPath: string): void
    symlinkSync(existingPath: string, newPath: string): void
    unlinkSync(filePath: string): void
    mkdirSync(dirPath: string, recursive: boolean, flag: number): void
    mkdtempSync(dirPath: string, prefix: string): string
    renameSync(oldPath: string, newPath: string): void
    rmSync(dirPath: string, recursive: boolean, force: boolean): void
}

const modFS = progpGetModule<ModFS>("nodejsModFS")!;

export const existsSync = modFS.existsSync;
export const chmodSync = modFS.chmodSync;
export const chownSync = modFS.chownSync;
export const truncateSync = modFS.truncateSync;
export const copyFileSync = modFS.copyFileSync;
export const linkSync = modFS.linkSync;
export const symlinkSync = modFS.symlinkSync;
export const unlinkSync = modFS.unlinkSync;
export const mkdtempSync = modFS.mkdtempSync;
export const renameSync = modFS.renameSync;

interface StatSyncOptions {
    throwIfNoEntry?: boolean
}

export function statSync(path: string, options: StatSyncOptions): any {
    let throwIfNoEntry = true;
    if (options && !options.throwIfNoEntry) throwIfNoEntry = false;
    return modFS.statSync(path, throwIfNoEntry)
}

export function readlinkSync(path: string, options?: ReadFileOptions): string{
    let encoding = "";

    if (options) {
        if (options.encoding) encoding = options.encoding;
    }

    if (encoding=="utf8") {
        return modFS.readFileUtf8Sync(path);
    } else {
        let bytes = modFS.readFileBytesSync(path);
        // TODO: must create a buffer
        throw "not implemented yet"
    }
}

export function mkdirSync(dirPath: string, options?: MkDirSyncOptions) {
    let recursive = false;
    let flag = 0o777;

    if (options) {
        if (options.recursive!==undefined) recursive = !!options.recursive

        if (options.mode!==undefined) {
            if (typeof(options.mode)=="string") {
                flag = fsOctalStringToInt(options.mode)
            }
        }
    }

    modFS.mkdirSync(dirPath, recursive, flag);
}

export function rmSync(dirPath: string, options?: RmSyncOptions) {
    let recursive = false, force = false;

    if (options) {
        if (options.recursive!==undefined) recursive = !!options.recursive
        if (options.force!==undefined) force = !!options.force
    }

    modFS.rmSync(dirPath, recursive, force);
}

export const rmdirSync = rmSync;

//region Const & Interfaces

interface ReadFileOptions {
    encoding?: string|null
    flag?: string
}

interface MkDirSyncOptions {
    recursive?: boolean
    mode?: string|number
}

interface RmSyncOptions {
    recursive?: boolean
    force?: boolean
}

export const constants = {
    UV_FS_SYMLINK_DIR: 1,
    UV_FS_SYMLINK_JUNCTION: 2,
    O_RDONLY: 0,
    O_WRONLY: 1,
    O_RDWR: 2,
    UV_DIRENT_UNKNOWN: 0,
    UV_DIRENT_FILE: 1,
    UV_DIRENT_DIR: 2,
    UV_DIRENT_LINK: 3,
    UV_DIRENT_FIFO: 4,
    UV_DIRENT_SOCKET: 5,
    UV_DIRENT_CHAR: 6,
    UV_DIRENT_BLOCK: 7,
    EXTENSIONLESS_FORMAT_JAVASCRIPT: 0,
    EXTENSIONLESS_FORMAT_WASM: 1,
    S_IFMT: 61440,
    S_IFREG: 32768,
    S_IFDIR: 16384,
    S_IFCHR: 8192,
    S_IFBLK: 24576,
    S_IFIFO: 4096,
    S_IFLNK: 40960,
    S_IFSOCK: 49152,
    O_CREAT: 512,
    O_EXCL: 2048,
    UV_FS_O_FILEMAP: 0,
    O_NOCTTY: 131072,
    O_TRUNC: 1024,
    O_APPEND: 8,
    O_DIRECTORY: 1048576,
    O_NOFOLLOW: 256,
    O_SYNC: 128,
    O_DSYNC: 4194304,
    O_SYMLINK: 2097152,
    O_NONBLOCK: 4,
    S_IRWXU: 448,
    S_IRUSR: 256,
    S_IWUSR: 128,
    S_IXUSR: 64,
    S_IRWXG: 56,
    S_IRGRP: 32,
    S_IWGRP: 16,
    S_IXGRP: 8,
    S_IRWXO: 7,
    S_IROTH: 4,
    S_IWOTH: 2,
    S_IXOTH: 1,
    F_OK: 0,
    R_OK: 4,
    W_OK: 2,
    X_OK: 1,
    UV_FS_COPYFILE_EXCL: 1,
    COPYFILE_EXCL: 1,
    UV_FS_COPYFILE_FICLONE: 2,
    COPYFILE_FICLONE: 2,
    UV_FS_COPYFILE_FICLONE_FORCE: 4,
    COPYFILE_FICLONE_FORCE: 4
}

//endregion

//region Tools functions

function fsOctalStringToInt(mode: string): number {
    return parseInt(mode, 8);
}

//endregion

export default {
    existsSync: existsSync,
    statSync: statSync,
    chmodSync: chmodSync,
    chownSync: chownSync,
    truncateSync: truncateSync,
    readlinkSync: readlinkSync,
    copyFileSync: copyFileSync,
    linkSync: linkSync,
    symlinkSync: symlinkSync,
    unlinkSync: unlinkSync,
    mkdtempSync: mkdtempSync,
    renameSync: renameSync,
    rmSync: rmSync,
    rmdirSync: rmdirSync,

    constants: constants,
}