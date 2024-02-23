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

package modReact

import (
	"embed"
	"github.com/progpjs/progpjs/v2"
)

//go:embed embed/*
var gEmbedFS embed.FS

func registerEmbeddedModule(innerPath string, alias ...string) {
	ebPath := progpjs.ReturnEmbeddedTypescriptModule(gEmbedFS, "embed/"+innerPath)

	for _, e := range alias {
		progpjs.AddJavascriptModuleProvider(e, ebPath)
	}
}

func InstallProgpJsModule() {
	// Here "@progp/react" is used when adding an alias
	// while "react" is used when react isn't found anywhere.
	//
	registerEmbeddedModule("jsMods/@progp/react/index.ts", "@progp/react", "react")
}
