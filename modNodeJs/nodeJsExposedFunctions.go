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
	"encoding/json"
	"github.com/progpjs/libProgpScripts"
	"github.com/progpjs/progpAPI"
	"os"
	"runtime"
)

func registerExportedFunctions() {
	rg := libProgpScripts.GetFunctionRegistry()
	myMod := rg.UseGoNamespace("github.com/progpjs/modules/modNodeJs")
	group := myMod.UseCustomGroup("nodejsModProcess")

	group.AddFunction("cwd", "JsCwd", JsCwd)
	group.AddFunction("env", "JsEnv", JsEnv)
	group.AddFunction("arch", "JsArch", JsArch)
	group.AddFunction("platform", "JsPlatform", JsPlatform)
	group.AddFunction("argv", "JsArgV", JsArgV)
	group.AddFunction("exit", "JsExit", JsExit)
	group.AddFunction("pid", "JsPID", JsPID)
	group.AddFunction("ppid", "JsPpID", JsPpID)
}

func JsCwd() string {
	cwd, _ := os.Getwd()
	return cwd
}

func JsEnv() progpAPI.StringBuffer {
	res := os.Environ()
	b, _ := json.Marshal(res)
	return b
}

func JsArch() string {
	// Apple MAC: arm64
	return runtime.GOARCH
}

func JsPlatform() string {
	return runtime.GOOS
}

func JsArgV() []string {
	return os.Args
}

func JsExit(code int) {
	os.Exit(code)
}

func JsPID() int {
	return os.Getpid()
}

func JsPpID() int {
	return os.Getppid()
}
