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

// https://nodejs.org/api/stream.html
// https://developer.mozilla.org/en-US/docs/Web/JavaScript/Guide/Typed_arrays
// https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/DataView

class Buffer extends Uint8Array {
    write(value: any, offset?: number, length?: number, _?: any): number {
        let ab = toArrayBuffer(value);

        if (typeof (offset)=="object") {
            offset = 0;
            length = ab.byteLength;
        } else if (typeof length=="object") {
            length = ab.byteLength;
        }

        return writeArrayBuffer(this, ab, offset!, length!);
    }
}

export function alloc(size: number): Buffer {
    return new Buffer(size);
}

export function byteLength(buffer: Buffer): number {
    return buffer.byteLength;
}

export function from(value: any, offset?: number, length?: number): Buffer {
    let ab = toArrayBuffer(value);

    if (offset===undefined) offset = 0;
    if (length===undefined) length = ab.byteLength;
    let buffer = alloc(offset+length);
    writeArrayBuffer(buffer, ab, offset, length);

    return buffer;
}

//region Tools

function toArrayBuffer(v: any): ArrayBuffer {
    if (v instanceof ArrayBuffer) return v;
    if (v instanceof Uint8Array) return v.buffer;
    if (v.substring!==undefined) return progpStringToBuffer(v);

    throw "not implemented";
}

function writeArrayBuffer(buffer: Buffer, ab: ArrayBuffer, offset: number, length: number): number {
    let ui8 = new Uint8Array(ab, 0, length);
    buffer.set(ui8, offset);
    return length;
}

//endregion

export default {
    alloc: alloc
};