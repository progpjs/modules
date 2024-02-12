// https://nodejs.org/api/process.html

interface ModProcess {
    cwd(): string
}

const modProcess = progpGetModule<ModProcess>("nodejsModProcess")!;

export const cwd = modProcess.cwd;

export default {
    cwd: cwd
}