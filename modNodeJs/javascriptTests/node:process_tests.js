const process = require('node:process');
const test = require("node:test");
const assert = require("node:assert");

test("NodeJS 'process'", () => {

    // Supported but commented in order to avoid
    // stopping the current test.
    //
    //process.exit(5);
    //process.kill(process.pid, 0);

    console.log("process.pid: ", process.pid);
    console.log("process.ppid: ", process.ppid);
    console.log("process.getuid: ", process.getuid());

    process.chdir("/"); // Throw error if directory don't exist.
    console.log(`process.cwd(): ${process.cwd()}`);

    console.log("process.arch: ", process.arch);
    console.log("process.platform: ", process.platform);
    console.log("process.argv: ", process.argv);
    console.log("process.argv0: ", process.argv0);
    console.log("process.execArgv: ", process.execArgv);
    console.log("process.execPath: ", process.execPath);
    console.log("process.env.PATH: ", process.env.PATH);

    process.nextTick(() => {
        console.log("on next tick !!!")
    });
});
