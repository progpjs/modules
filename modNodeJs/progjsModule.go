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

package modNodeJs

import (
	"embed"
	"github.com/progpjs/progpjs"
)

//go:embed embed/*
var gEmbedFS embed.FS

func registerEmbeddedModule(innerPath string, alias ...string) {
	provider := progpjs.ReturnEmbeddedTypescriptModule(gEmbedFS, "embed/"+innerPath)

	for _, e := range alias {
		progpjs.AddJavascriptModuleProvider(e, provider)
	}
}

func InstallModule() {
	registerExportedFunctions()

	registerEmbeddedModule("jsMods/assert.ts", "assert", "node:assert")
	registerEmbeddedModule("jsMods/test.ts", "test", "node:test")

	registerEmbeddedModule("jsMods/fs.ts", "fs", "node:fs")
	registerEmbeddedModule("jsMods/os.ts", "os", "node:os")
	registerEmbeddedModule("jsMods/path.ts", "path", "node:path")
	registerEmbeddedModule("jsMods/process.ts", "process", "node:process")
	registerEmbeddedModule("jsMods/stream.ts", "stream", "node:stream")
	registerEmbeddedModule("jsMods/buffer.ts", "buffer", "node:buffer")
	registerEmbeddedModule("jsMods/timers.ts", "timers", "node:timers")
	registerEmbeddedModule("jsMods/url.ts", "url", "node:url")
}
