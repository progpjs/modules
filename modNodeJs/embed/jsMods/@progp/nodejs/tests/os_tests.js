import test from "node:test";
import os from "node:os";

test("NodeJS 'os'", () => {
    console.log("os.hostname: ", os.hostname());
    console.log("os.homedir: ", os.homedir());
    console.log("os.tmpdir: ", os.tmpdir());
    console.log("os.arch: ", os.arch());
    console.log("os.machine: ", os.machine());
    console.log("os.platform: ", os.platform());
});