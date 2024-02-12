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
	"syscall"
)

func registerExportedFunctions() {
	rg := libProgpScripts.GetFunctionRegistry()
	myMod := rg.UseGoNamespace("github.com/progpjs/modules/modNodeJs")

	modProcess := myMod.UseCustomGroup("nodejsModProcess")
	modProcess.AddFunction("kill", "JsKill", JsKill)
	modProcess.AddFunction("cwd", "JsCwd", JsCwd)
	modProcess.AddFunction("env", "JsEnv", JsEnv)
	modProcess.AddFunction("arch", "JsArch", JsArch)
	modProcess.AddFunction("platform", "JsPlatform", JsPlatform)
	modProcess.AddFunction("argv", "JsArgV", JsArgV)
	modProcess.AddFunction("exit", "JsExit", JsExit)
	modProcess.AddFunction("pid", "JsPID", JsPID)
	modProcess.AddFunction("ppid", "JsPpID", JsPpID)
	modProcess.AddFunction("chdir", "JsChDir", JsChDir)
	modProcess.AddFunction("getuid", "JsGetUid", JsGetUid)
	modProcess.AddAsyncFunction("nextTick", "JsNextTickAsync", JsNextTickAsync)

	modOS := myMod.UseCustomGroup("nodejsModOS")
	modOS.AddFunction("homeDir", "JsHomeDir", JsHomeDir)
	modOS.AddFunction("hostName", "JsHostName", JsHostName)
	modOS.AddFunction("tempDir", "JsTempDir", JsTempDir)
}

//region node:process	(nodejsModProcess)

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

func JsChDir(dir string) error {
	return os.Chdir(dir)
}

func JsGetUid() int {
	return os.Getuid()
}

func JsNextTickAsync(fct progpAPI.ScriptFunction) {
	progpAPI.SafeGoRoutine(func() {
		fct.CallWithUndefined()
	})
}

func JsKill(pid int, signal int) error {
	err := syscall.Kill(pid, syscall.Signal(signal))

	if signal == 0 {
		// Don't throw error fi signal is 0
		// which allows testing if process exists.
		// It's a node.js special case.
		//
		return nil
	}

	return err
}

//endregion

//region node:os (nodejsModOS)

func JsHomeDir() (string, error) {
	dirname, err := os.UserHomeDir()
	return dirname, err
}

func JsHostName() (string, error) {
	name, err := os.Hostname()
	return name, err
}

func JsTempDir() string {
	return os.TempDir()
}

//endregion
