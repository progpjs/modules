const process = require('node:process');
const test = require("node:test");
const assert = require("node:assert");

test("NodeJS 'process.cwd()'", () => {
    let cwd = process.cwd();

    //process.exit(5);

    console.log("process.pid: ", process.pid);
    console.log("process.ppid: ", process.ppid);
    console.log(`process.cwd(): ${process.cwd()}`);
    console.log("process.arch: ", process.arch);
    console.log("process.platform: ", process.platform);
    console.log("process.argv: ", process.argv);
    console.log("process.argv0: ", process.argv0);
    console.log("process.execArgv: ", process.execArgv);
    console.log("process.execPath: ", process.execPath);
    console.log("process.env.PATH: ", process.env.PATH);
});
