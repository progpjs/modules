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
	"github.com/progpjs/progpjs/v2"
	"strings"
)

//go:embed embed/*
var gEmbedFS embed.FS

func registerEmbeddedModule(innerPath string, alias ...string) {
	if strings.HasPrefix("node:", innerPath) {
		innerPath = innerPath[5:]
	}

	ebPath := progpjs.ReturnEmbeddedTypescriptModule(gEmbedFS, "embed/"+innerPath)

	for _, e := range alias {
		progpjs.AddJavascriptModuleProvider(e, ebPath)
	}
}

func InstallProgpJsModule() {
	registerExportedFunctions()

	registerEmbeddedModule("jsMods/@progp/nodejs/assert.ts", "assert", "node:assert")
	registerEmbeddedModule("jsMods/@progp/nodejs/test.ts", "test", "node:test")

	registerEmbeddedModule("jsMods/@progp/nodejs/fs.ts", "fs", "node:fs")
	registerEmbeddedModule("jsMods/@progp/nodejs/os.ts", "os", "node:os")
	registerEmbeddedModule("jsMods/@progp/nodejs/path.ts", "path", "node:path")
	registerEmbeddedModule("jsMods/@progp/nodejs/process.ts", "process", "node:process")
	registerEmbeddedModule("jsMods/@progp/nodejs/stream.ts", "stream", "node:stream")
	registerEmbeddedModule("jsMods/@progp/nodejs/buffer.ts", "buffer", "node:buffer")
	registerEmbeddedModule("jsMods/@progp/nodejs/timers.ts", "timers", "node:timers")
	registerEmbeddedModule("jsMods/@progp/nodejs/url.ts", "url", "node:url")
}
