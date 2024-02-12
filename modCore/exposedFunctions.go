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

package modCore

import (
	"github.com/progpjs/libProgpScripts"
	"github.com/progpjs/progpAPI"
)

func registerExportedFunctions() {
	rg := libProgpScripts.GetFunctionRegistry()
	myMod := rg.UseGoNamespace("github.com/progpjs/modules/modCore")
	group := myMod.UseGroupGlobal()

	group.AddAsyncFunction("progpCallAfterMs", "JsProgpCallAfterMsAsync", JsProgpCallAfterMsAsync)
}

func JsProgpCallAfterMsAsync(timeInMs int, callback progpAPI.ScriptFunction) {
	progpAPI.SafeGoRoutine(func() {
		progpAPI.PauseMs(timeInMs)
		callback.CallWithUndefined()
	})
}
