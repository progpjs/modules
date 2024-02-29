const test = require('node:test');
const os = require("node:os");
const assert = require("node:assert");
const path = require("node:path");
const fs = require("node:fs");

test("NodeJS 'fs:sync'", () => {
    let homeDir = os.homedir();

    assert.strictEqual(fs.existsSync(homeDir), true, "fs.existsSync - 1")
    assert.strictEqual(fs.existsSync(path.join(homeDir, "dontExists")), false, "fs.existsSync - 2")

    console.log("statSync: ", fs.statSync(path.join(homeDir, ".profile")));
});