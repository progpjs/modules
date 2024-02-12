// https://nodejs.org/api/path.html

import process from "node:process";

const sep = "/";

export function basename(path: any, suffix?: string): string {
    if (!path) return "";

    let idx = path.lastIndexOf('/');
    if (idx!=-1) path = path.substring(idx + 1);

    if (suffix!==undefined) {
        if (path.endsWith(suffix)) {
            return path.substring(0, path.length - suffix.length);
        }
    }

    return path;
}

// https://nodejs.org/api/path.html#pathdelimiter
export const delimiter = ":";

// https://nodejs.org/api/path.html#pathdirnamepath
export function dirname(path: string): string {
    if (!path) return ".";

    if (path[path.length-1]=="/") {
        path = path.substring(0, path.length-1);
        if (!path) return "/"
    }

    let idx = path.lastIndexOf('/');
    if (idx==-1) return ".";

    return path.substring(0, idx);
}

// https://nodejs.org/api/path.html#pathextnamepath
export function extname(path: string) {
    if (!path) return "";

    let idx = path.lastIndexOf('.');
    if (idx==-1) return "";
    if (idx==0) return "";

    return path.substring(idx);
}

// https://nodejs.org/api/path.html#pathjoinpaths
//
export function join(...paths: string[]) {
    if (!paths) return ".";
    let size = paths.length;
    if (!size) return ".";

    let res = paths[0];
    let i = 1;

    if (res=="") {
        do {
            res = paths[i++];
        } while ((res=="")&&(i<size));
    }

    if (!res) return ".";
    let endsWithSep = res[res.length-1] == "/";

    for (;i<size;i++) {
        let p = paths[i];
        if (p=="") continue;

        if (endsWithSep) {
            if (p[0] == "/") res += p.substring(1);
            else res += p;
        } else {
            if (p[0] == "/") res += p;
            else res += "/" + p;
        }

        endsWithSep = p[p.length-1] == "/"
    }

    // Cas: "//a" --> "/a".
    do {
        if (res.length<=1) return res;

        if (res[0] == "/") {
            if (res[1] == "/") {
                res = res.substring(1);
            } else {
                break;
            }
        } else {
            break;
        }
    } while (true);

    // Cas: "a//" --> "a/".
    //
    if (endsWithSep) {
        i = res.length-2;

        while (res[i]=="/") {
            res = res.substring(0, i+1);
            i--;
        }
    }

    if (!res) return ".";
    return res;
}

interface PathObject {
    dir?: string;
    root?: string;
    base?: string;
    name?: string;
    ext?: string;
}

// https://nodejs.org/api/path.html#pathformatpathobject
//
export function format(pathObject: PathObject) {
    let p = "";

    if (pathObject.base) {
        p += pathObject.base;
    } else {
        if (pathObject.name) p += pathObject.name;

        if (pathObject.ext) {
            if (pathObject.ext[0]!=".") p += "."
            p += pathObject.ext;
        }
    }

    if (pathObject.dir) {
        p = join(pathObject.dir, p)
    } else if (pathObject.root) {
        // root is for windows, it's for exemple "c:\" or "d:\".
        // Here windows isn't supported but the behaviors is.
        //
        p = join(pathObject.root, p)
    }

    return p;
}

// https://nodejs.org/api/path.html#pathparsepath
export function parse(path: string): PathObject {
    let res: PathObject = {
        root: "", base: "", dir: "", ext: "", name: ""
    };

    if (!path) return res;

    let idx = path.lastIndexOf("/");
    if (idx===-1) {
        res.dir = "";
        res.base = path;
    } else {
        if (idx===0) res.dir = "/";
        else res.dir = path.substring(0, idx);
        res.base = path.substring(idx+1);
    }

    idx = res.base.lastIndexOf(".");
    if (idx!==-1) {
        res.name = res.base.substring(0, idx)
        res.ext = res.base.substring(idx);
    }

    if (res.dir) {
        idx = res.dir.indexOf("/");
        res.root = res.dir.substring(0, idx+1);
    }

    return res;
}

// https://nodejs.org/api/path.html#pathtonamespacedpathpath
export function toNamespacedPath(path: string): string {
    // Is only relevant on windows.
    // With posix system returns the path as is.
    return path;
}

// https://nodejs.org/api/path.html#pathresolvepaths
export function resolve(...paths: string[]) {
    let segments: string[] = [];
    let isRel = true;
    let isFirst = true;

    for (let p of paths) {
        if (!p) continue;

        if (p[0]=='/') {
            segments = [];
            isRel = false;
            isFirst = true;
        }

        let parts = p.split("/");

        for (let part of parts) {
            if (!isFirst && !part) continue;
            isFirst = false;

            if (part==".") {
                continue;
            }

            if (part=="..") {
                if (segments.length!==0) segments.pop();
                continue;
            }

            segments.push(part);
        }
    }

    let res = segments.join("/");
    if (isRel) res = join(process.cwd(), res);

    return res;
}

export default {
    sep: sep,
    delimiter: delimiter,

    basename: basename,
    extname: extname,
    join: join,
    format: format,
    parse: parse,
    dirname: dirname,
}