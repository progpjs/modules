function test(testName: string, testFunction: Function) {
    console.log("Executing test: ", testName);
    let res;

    try {
        res = testFunction();
    }
    catch (e) {
        // TODO
        return;
    }

    if (res instanceof Promise) {
        res.then(value => {
            // TODO
        });

        res.catch(err => {
            // TODO
        });
    }
}

interface ErrorInfos {
    title: string;
    userMessage?: string|Error;

    hidedExpectedValue?: boolean;
    expectValue?: any;
    foundValue?: any;
}

function error(infos: ErrorInfos) {
    Assert._errorCount++;

    console.log("[ASSERTION ERROR] - " + infos.title);
    if (infos.userMessage) console.log("    |- Test: ", infos.userMessage);

    if (!infos.hidedExpectedValue) {
        let v = infos.expectValue;

        if ((v!==undefined)&&(v!==null)) {
            if (v.toString) v = v.toString();
            console.log("    |- Expected: [", v, "]");
        } else {
            if (v===undefined) console.log("    |- Expected undefined");
            else console.log("    |- Expected null");
        }

        v = infos.foundValue;

        if ((v!==undefined)&&(v!==null)) {
            if (v.toString) v = v.toString();
            console.log("    |- Found: [", v, "]");
        } else {
            if (v===undefined) console.log("    |- Found undefined");
            else console.log("    |- Found null");
        }
    }
}

const Assert = {
    // Used by the test function to count errors.
    _errorCount: 0,

    strictEqual: function(actual: any, expected: any, message?: string|Error) {
        if (actual===expected) return;
        error({title: "Expected values to be strictly equal", userMessage: message, foundValue: actual, expectValue: expected});
    },

    equal: function(actual: any, expected: any, message?: string|Error) {
        if (actual==expected) return;
        error({title: "Expected values to be equal", userMessage: message, foundValue: actual, expectValue: expected});
    },

    throws: function(f: Function, message?: string|Error) {
        let hasError = false;
        try { f() } catch (e) { hasError = true; }

        if (!hasError) error({title: "Expected throwing an error", userMessage: message, hidedExpectedValue: true});
    }
}

// @ts-ignore
module.exports = Assert;
