// https://nodejs.org/api/process.html

interface ModProcess {
    cwd(): string
    env(): string
    argv(): string[]
    arch(): string
    platform(): string
    exit(code: number): void
    pid(): number
    ppid(): number
}

const modProcess = progpGetModule<ModProcess>("nodejsModProcess")!;

export const cwd = modProcess.cwd;

export const env = (() => {
    let res: any = {};
    let raw = JSON.parse(modProcess.env())

    for (let e of raw) {
        let idx = e.indexOf("=");
        res[e.substring(0, idx)] = e.substring(idx+1);
    }

    return res;
})();

export function exit(code: number) {
    if (!code) code = 0;
    modProcess.exit(code);
}

export const argv = modProcess.argv();
export const argv0 = argv[0];
export const execPath = argv0;
export const arch = modProcess.arch()
export const platform = modProcess.platform()
export const pid = modProcess.pid()
export const ppid = modProcess.ppid()

export const execArgv = [];

export default {
    cwd: cwd,
    env: env,
    arch: arch,
    platform: platform,
    argv: argv,
    argv0: argv0,
    execArgv: execArgv,
    execPath: execPath,
    exit: exit,
    pid: pid,
    ppid: ppid,
}