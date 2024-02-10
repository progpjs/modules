// @ts-ignore
const assert = require("assert");

// https://nodejs.org/api/test.html#test-runner

function defaultFunction(testName: string, testFunction: Function) {
    function endTest() {
        if (assert._errorCount==0) {
            console.log("✔ Test success: " + testName)
        }
    }

    let res;

    assert._errorCount = 0;

    try {
        res = testFunction();
    }
    catch (e) {
        let asString: string;
        if (e instanceof Error) asString = e.toString(); else asString = e as string;

        console.error("Unexpected error :", asString);
        assert._errorCount++;
    }

    if (res instanceof Promise) {
        res.then(value => {
            endTest();
        });

        res.catch(err => {
            assert._errorCount++;
            console.error("Unexpected error :", err);
            endTest();
        });
    } else {
        endTest();
    }
}

// Allows doing:
//      const test = require('node:test');
// or
//      import test from "node:test";
//
// @ts-ignore
module.exports = defaultFunction;
