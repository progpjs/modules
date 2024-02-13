const test = require('node:test');

test("NodeJS 'path.others'", () => {
    debugger;

    let buf = Buffer.alloc(100);
    console.log(buf.write("AA", 3, 1));
    console.log(buf.write("BB", 4, 1));
    console.log(buf);

    buf = Buffer.from("test")
    console.log(buf);
    console.log(buf.byteLength)
});