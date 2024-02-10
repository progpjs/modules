// https://nodejs.org/api/path.html

// https://nodejs.org/api/path.html#pathbasenamepath-suffix
//
// @ts-ignore
const pProto = "".prototype;

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

interface PathObject {
    dir?: string;
    root?: string;
    base?: string;
    name?: string;
    ext?: string;
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


// https://nodejs.org/api/path.html#pathformatpathobject
//
export function format(pathObject: PathObject) {

}

export default {
    basename: basename,
    delimiter: delimiter,
    extname: extname,
    join: join,
    format: format,
    dirname: dirname
}