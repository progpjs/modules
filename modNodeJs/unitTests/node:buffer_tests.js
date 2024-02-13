const test = require('node:test');

test("NodeJS 'path.others'", () => {
    let buf = Buffer.alloc(100);
    //console.log(buf.write("AA", 3, 1));
    //console.log(buf.write("BB", 4, 1));
    //console.log(buf.toJSON());

    buf = Buffer.from("test")
    //console.log(buf.toJSON());
    //console.log(buf.byteLength)

    buf = Buffer.concat([Buffer.from("AAA"), Buffer.from("BBB")], 10);
    //console.log(buf.toJSON())

    console.log(buf.slice(2, 4).toJSON());
    buf.fill(5)
    console.log(buf);
});