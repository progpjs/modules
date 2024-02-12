const test = require('node:test');
const assert = require("node:assert");
const path = require("node:path");

test("NodeJS 'path.others'", () => {
    assert.strictEqual(path.delimiter, ":");
});

test("NodeJS 'path.basename'", () => {
    assert.strictEqual(path.basename("f1.txt"), "f1.txt");
    assert.strictEqual(path.basename("/path/f1.txt"), "f1.txt");
    assert.strictEqual(path.basename(""), "");
    assert.strictEqual(path.basename("f1.txt.toRemove", ".toRemove"), "f1.txt");
});

test("NodeJS 'path.dirname'", () => {
    assert.strictEqual(path.dirname(""), ".");
    assert.strictEqual(path.dirname("noDir"), ".");
    assert.strictEqual(path.dirname("my/dir/f1.txt"), "my/dir");
    assert.strictEqual(path.dirname("./my/dir/f1.txt"), "./my/dir");
    assert.strictEqual(path.dirname("/root/dir/f1.txt"), "/root/dir");
    assert.strictEqual(path.dirname("/root/dir/f1.txt"), "/root/dir");

    assert.strictEqual(path.dirname("//root/dir/f1.txt"), "//root/dir");

    assert.strictEqual(path.dirname("//root/dir//f1.txt"), "//root/dir/");
    assert.strictEqual(path.dirname("//root/dir///f1.txt"), "//root/dir//");

    assert.strictEqual(path.dirname("//rel/dir/../f1.txt"), "//rel/dir/..");
    assert.strictEqual(path.dirname("/rel/dir/.."), "/rel/dir");
    assert.strictEqual(path.dirname("/rel/dir/"), "/rel");

    assert.strictEqual(path.dirname("/"), "/");
});

test("NodeJS 'path.extname'", () => {
    assert.strictEqual(path.extname('index.html'), ".html");
    assert.strictEqual(path.extname('index.coffee.md'), ".md");
    assert.strictEqual(path.extname('index.'), ".");
    assert.strictEqual(path.extname('index'), "");
    assert.strictEqual(path.extname('.index'), "");
    assert.strictEqual(path.extname('.index.md'), ".md");
});

test("NodeJS 'path.join'", () => {
    assert.strictEqual(path.join(""), ".",);
    assert.strictEqual(path.join(), ".");
    assert.strictEqual(path.join("/"), "/");

    assert.strictEqual(path.join("a", "b"), "a/b");
    assert.strictEqual(path.join("a", "b/"), "a/b/");
    assert.strictEqual(path.join("/a", "b/"), "/a/b/");

    assert.strictEqual(path.join("/a", "", "b/"), "/a/b/");
    assert.strictEqual(path.join("/a", "", "/b/"), "/a/b/");

    assert.strictEqual(path.join("", "b/"), "b/");
    assert.strictEqual(path.join("", "", "b/"), "b/");
    assert.strictEqual(path.join("", "", "b/", ""), "b/");
    assert.strictEqual(path.join("", "", "b", ""), "b");

    assert.strictEqual(path.join("a", "b//"), "a/b/");
    assert.strictEqual(path.join("a", "b///"), "a/b/");

    assert.strictEqual(path.join("/", "///"), "/");
    assert.strictEqual(path.join("/", "//a"), "/a");
    assert.strictEqual(path.join("/", "//", "a"), "/a");
});

test("NodeJS 'path.format'", () => {
    let res;

    res = path.format({root: '/ignored', dir: '/home/user/dir', base: 'file.txt'});
    assert.strictEqual(res, "/home/user/dir/file.txt", "path.format/2");

    res = path.format({root: '/', base: 'file.txt', ext: 'ignored'});
    assert.strictEqual(res, "/file.txt", "path.format/2");

    res = path.format({root: '/', name: 'file', ext: '.txt'});
    assert.strictEqual(res, "/file.txt","path.format/3");

    res = path.format({root: '/', name: 'file', ext: 'txt'});
    assert.strictEqual(res, "/file.txt", "path.format/4");
});