/*
 * (C) Copyright 2024 Johan Michel PIQUET, France (https://johanpiquet.fr/).
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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

// @ts-ignore
module.exports = defaultFunction;
