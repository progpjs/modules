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
	"errors"
	"github.com/progpjs/progpAPI/v2"
	"github.com/progpjs/progpjs/v2"
	"os"
	"path"
)

func registerExportedFunctions() {
	rg := progpjs.GetFunctionRegistry()
	myMod := rg.UseGoNamespace("github.com/progpjs/modules/v2/modCore")
	group := myMod.UseGroupGlobal()

	group.AddAsyncFunction("progpCallAfterMs", "JsProgpCallAfterMsAsync", JsProgpCallAfterMsAsync)
	group.AddFunction("progpDispose", "JsProgpDispose", JsProgpDispose)
	group.AddFunction("progpAutoDispose", "JsProgpAutoDispose", JsProgpAutoDispose)
	group.AddAsyncFunction("progpRunScript", "JsProgpRunScriptAsync", JsProgpRunScriptAsync)
	group.AddFunction("progpSignal", "JsProgpSignal", JsProgpSignal)
	group.AddFunction("progpReturnString", "JsProgpReturnString", JsProgpReturnString)
	group.AddFunction("progpReturnVoid", "JsProgpReturnVoid", JsProgpReturnVoid)
	group.AddFunction("progpReturnError", "JsProgpReturnError", JsProgpReturnError)
}

type ProgpReturnErrorAction interface {
	OnReturnErrorAction(err string) error
}

type ProgpReturnVoidAction interface {
	OnReturnVoidAction() error
}

type ProgpReturnStringAction interface {
	OnReturnStringAction(value string) error
}

func JsProgpSignal(rc *progpAPI.SharedResourceContainer, signal string, data string) error {
	return progpjs.EmitProgpSignal(rc.GetScriptContext(), signal, data)
}

func JsProgpReturnString(res *progpAPI.SharedResource, value string) error {
	if action, ok := res.Value.(ProgpReturnStringAction); ok {
		return action.OnReturnStringAction(value)
	}

	return errors.New("invalid return call")
}

func JsProgpReturnError(res *progpAPI.SharedResource, error string) error {
	if action, ok := res.Value.(ProgpReturnErrorAction); ok {
		return action.OnReturnErrorAction(error)
	}

	return errors.New("invalid return call")
}

func JsProgpReturnVoid(res *progpAPI.SharedResource) error {
	if action, ok := res.Value.(ProgpReturnVoidAction); ok {
		return action.OnReturnVoidAction()
	}

	return errors.New("invalid return call")
}

func JsProgpRunScriptAsync(rc *progpAPI.SharedResourceContainer, scriptFilePath string, securityGroup string, callback progpAPI.JsFunction) {
	if !path.IsAbs(scriptFilePath) {
		cwd, _ := os.Getwd()
		scriptFilePath = path.Join(cwd, scriptFilePath)
	}

	ctx := rc.GetScriptContext().GetScriptEngine().CreateNewScriptContext(securityGroup, false)

	progpAPI.SafeGoRoutine(func() {
		err := ctx.ExecuteChildScriptFile(scriptFilePath)

		if err != nil {
			callback.CallWithError(err)
		} else {
			callback.CallWithUndefined()
		}
	})
}

func JsProgpAutoDispose(rc *progpAPI.SharedResourceContainer, f progpAPI.JsFunction) {
	// Enable the auto-disposing mechanism for this function.
	f.EnabledResourcesAutoDisposing(rc)

	// Do the call himself.
	f.CallWithUndefined()
}

func JsProgpDispose(res progpAPI.SharedResource) {
	res.Dispose()
}

func JsProgpCallAfterMsAsync(timeInMs int, callback progpAPI.JsFunction) {
	progpAPI.SafeGoRoutine(func() {
		progpAPI.PauseMs(timeInMs)
		callback.CallWithUndefined()
	})
}
