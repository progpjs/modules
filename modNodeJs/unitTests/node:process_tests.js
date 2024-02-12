const { cwd } = require('node:process');
const test = require("node:test");
const assert = require("node:assert");

test("NodeJS 'process.cwd()'", () => {
    console.log(`Current directory: ${cwd()}`);
});

console.log("ici!")