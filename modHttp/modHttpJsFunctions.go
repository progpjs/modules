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

package modHttp

import (
	"bytes"
	"errors"
	"github.com/progpjs/libHttpServer"
	"github.com/progpjs/libHttpServer/libFastHttpImpl"
	"github.com/progpjs/progpAPI"
	"github.com/progpjs/progpScripts"
	"mime/multipart"
	"net/textproto"
	"sync"
)

var NoResponseSendError = errors.New("no response send")

type httpFormFile struct {
	Header textproto.MIMEHeader `json:"header"`
	Size   int64                `json:"size"`
	Id     int                  `json:"id"`
}

func registerExportedFunctions() {
	rg := progpScripts.GetFunctionRegistry()
	myMod := rg.UseGoNamespace("github.com/progpjs/modules/modHttp")
	group := myMod.UseCustomGroup("progpjsModHttp")

	group.AddFunction("startServer", "JsStartServer", JsStartServer)
	group.AddFunction("configureServer", "JsConfigureServer", JsConfigureServer)
	group.AddFunction("getHost", "JsGetHost", JsGetHost)

	group.AddFunction("VERB_withFunction", "JsVerbWithFunction", JsVerbWithFunction)

	group.AddFunction("returnString", "JsReturnString", JsReturnString)
	group.AddFunction("requestURI", "JsRequestURI", JsRequestURI)

	group.AddFunction("requestPath", "JsRequestPath", JsRequestPath)
	group.AddFunction("requestIP", "JsRequestIP", JsRequestIP)
	group.AddFunction("requestMethod", "JsRequestMethod", JsRequestMethod)
	group.AddFunction("requestHost", "JsRequestHost", JsRequestHost)
	group.AddFunction("requestQueryArgs", "JsRequestQueryArgs", JsRequestQueryArgs)
	group.AddFunction("requestPostArgs", "JsRequestPostArgs", JsRequestPostArgs)
	group.AddAsyncFunction("requestReadFormFile", "JsRequestReadFormFileAsync", JsRequestReadFormFileAsync)
	group.AddAsyncFunction("requestSaveFormFile", "JsRequestSaveFormFileAsync", JsRequestSaveFormFileAsync)

	group.AddFunction("requestWildcards", "JsRequestWildcards", JsRequestWildcards)
	group.AddFunction("requestRemainingSegments", "JsRequestRemainingSegments", JsRequestRemainingSegments)

	group.AddFunction("requestCookies", "JsRequestCookies", JsRequestCookies)
	group.AddFunction("requestCookie", "JsRequestCookie", JsRequestCookie)
	group.AddFunction("requestHeaders", "JsRequestHeaders", JsRequestHeaders)
	group.AddFunction("responseSetHeader", "JsRequestSetHeader", JsRequestSetHeader)
	group.AddFunction("responseSetCookie", "JsRequestSetCookie", JsRequestSetCookie)
}

// JsConfigureServer configure a server designed by his port.
// It does nothing if the server is already started, but returns false if the configuration can't be applied.
func JsConfigureServer(serverPort int, params libHttpServer.HttpServerStartParams) bool {
	server := libFastHttpImpl.GetFastHttpServer(serverPort)

	if server.IsStarted() {
		return false
	}

	server.SetStartServerParams(libHttpServer.HttpServerStartParams(params))
	return true
}

// JsGetHost returns an HttpHost object from a port and a hostname.
func JsGetHost(rc *progpAPI.SharedResourceContainer, serverPort int, hostName string) *progpAPI.SharedResource {
	server := libFastHttpImpl.GetFastHttpServer(serverPort)

	host := server.GetHost(hostName)
	return rc.NewSharedResource(host, nil)
}

// JsStartServer starts the server designed by his port.
// This server must have been configured before, otherwise it uses the default configuration.
func JsStartServer(rc *progpAPI.SharedResourceContainer, serverPort int) error {
	// Allows avoiding exiting the javascript VM.
	ctx := rc.GetScriptContext()
	ctx.IncreaseRefCount()

	server := libFastHttpImpl.GetFastHttpServer(serverPort)

	if server.IsStarted() {
		return nil
	}

	progpAPI.DeclareBackgroundTaskStarted()

	var mutex sync.Mutex
	mutex.Lock()

	var err error

	progpAPI.SafeGoRoutine(func() {
		mutex.Unlock()

		// Will block
		err = server.StartServer()
	})

	mutex.Lock()

	// This pause allows knowing if the server starts with an error.
	progpAPI.PauseMs(5)

	return err
}

// JsVerbWithFunction bind a GET/POST/... call to a function inside a context.
// This function is executed when the GET request match.
func JsVerbWithFunction(rc *progpAPI.SharedResourceContainer, resHost *progpAPI.SharedResource, verb string, requestPath string, callback progpAPI.ScriptFunction) error {
	host, ok := resHost.Value.(*libHttpServer.HttpHost)
	if !ok {
		return errors.New("invalid resource")
	}

	// Allows calling this function more than one time.
	callback.KeepAlive()

	host.VERB(verb, requestPath, func(call libHttpServer.HttpRequest) error {
		var mutex sync.Mutex
		call.SetUnlockMutex(&mutex)

		res := rc.NewSharedResource(call, nil)

		// Allows disposing before the host script exit.
		defer res.Dispose()

		mutex.Lock()
		callback.CallWithResource2(res)

		// Will block the call until a response is sent.
		mutex.Lock()

		if !call.IsBodySend() {
			return NoResponseSendError
		}

		return nil
	})

	return nil
}

// JsReturnString set the response to returns.
func JsReturnString(resHttpRequest *progpAPI.SharedResource, responseCode int, contentType string, responseText string) error {
	call, ok := resHttpRequest.Value.(libHttpServer.HttpRequest)
	if !ok {
		return errors.New("invalid resource")
	}

	call.SetContentType(contentType)

	// This response will unlock the caller mutex and the response will be sent.
	call.ReturnString(responseCode, responseText)
	return nil
}

func JsRequestURI(resHttpRequest *progpAPI.SharedResource) (error, string) {
	call, ok := resHttpRequest.Value.(libHttpServer.HttpRequest)
	if !ok {
		return errors.New("invalid resource"), ""
	}

	return nil, call.URI()
}

func JsRequestPath(resHttpRequest *progpAPI.SharedResource) (error, string) {
	call, ok := resHttpRequest.Value.(libHttpServer.HttpRequest)
	if !ok {
		return errors.New("invalid resource"), ""
	}

	return nil, call.Path()
}

func JsRequestIP(resHttpRequest *progpAPI.SharedResource) (error, string) {
	call, ok := resHttpRequest.Value.(libHttpServer.HttpRequest)
	if !ok {
		return errors.New("invalid resource"), ""
	}

	return nil, call.RemoteIP()
}

func JsRequestMethod(resHttpRequest *progpAPI.SharedResource) (error, string) {
	call, ok := resHttpRequest.Value.(libHttpServer.HttpRequest)
	if !ok {
		return errors.New("invalid resource"), ""
	}

	return nil, call.GetMethodName()
}

func JsRequestHost(resHttpRequest *progpAPI.SharedResource) (error, string) {
	call, ok := resHttpRequest.Value.(libHttpServer.HttpRequest)
	if !ok {
		return errors.New("invalid resource"), ""
	}

	return nil, call.GetHost().GetHostName()
}

func JsRequestQueryArgs(resHttpRequest *progpAPI.SharedResource) (error, map[string]string) {
	call, ok := resHttpRequest.Value.(libHttpServer.HttpRequest)
	if !ok {
		return errors.New("invalid resource"), nil
	}

	resMap := make(map[string]string)

	call.GetQueryArgs().VisitAll(func(key, value []byte) {
		resMap[string(key)] = string(value)
	})

	return nil, resMap
}

func JsRequestPostArgs(resHttpRequest *progpAPI.SharedResource) (error, map[string]any) {
	call, ok := resHttpRequest.Value.(libHttpServer.HttpRequest)
	if !ok {
		return errors.New("invalid resource"), nil
	}

	if !call.IsMultipartForm() {
		resMap := make(map[string]any)

		call.GetPostArgs().VisitAll(func(key, value []byte) {
			resMap[string(key)] = string(value)
		})

		return nil, resMap
	}

	form, err := call.GetMultipartForm()
	if err != nil {
		return err, nil
	}

	resMap := make(map[string]any)

	for k, v := range form.Values {
		if len(v) == 1 {
			resMap[k] = v[0]
		} else {
			resMap[k] = v
		}
	}

	if form.Files != nil {
		for formFieldName, files := range form.Files {
			var fileList []httpFormFile

			for id, entry := range files {
				file := httpFormFile{
					Header: entry.Header,
					Size:   entry.Size,
					Id:     id,
				}

				fileList = append(fileList, file)
			}

			if fileList != nil {
				// If it's a file, then "resMap" contains an array.
				// It must be kept as-is, to help javascript to detect.
				resMap[formFieldName] = fileList
			}
		}
	}

	return nil, resMap
}

func openRequestFormFile(resHttpRequest *progpAPI.SharedResource, fieldName string, fileOffset int, callback progpAPI.ScriptFunction) (multipart.File, bool) {
	call, ok := resHttpRequest.Value.(libHttpServer.HttpRequest)
	if !ok {
		callback.CallWithError(errors.New("invalid resource"))
		return nil, false
	}

	if !call.IsMultipartForm() {
		callback.CallWithError(errors.New("invalid field name"))
		return nil, false
	}

	mpf, err := call.GetMultipartForm()
	if err != nil {
		callback.CallWithError(err)
		return nil, false
	}

	if mpf.Files == nil {
		callback.CallWithError(errors.New("invalid field name"))
		return nil, false
	}

	files := mpf.Files[fieldName]
	if files == nil {
		callback.CallWithError(errors.New("invalid field name"))
		return nil, false
	}
	if (fileOffset < 0) || (fileOffset > len(files)) {
		callback.CallWithError(errors.New("invalid file id"))
		return nil, false
	}

	entry := files[fileOffset]

	file, err := entry.Open()
	if err != nil {
		callback.CallWithError(errors.New("can't read file"))
		return nil, false
	}

	return file, true
}

func JsRequestSaveFormFileAsync(resHttpRequest *progpAPI.SharedResource, fieldName string, fileOffset int, saveFilePath string, callback progpAPI.ScriptFunction) {
	file, ok := openRequestFormFile(resHttpRequest, fieldName, fileOffset, callback)
	if !ok {
		return
	}

	defer file.Close()

	err := libHttpServer.SaveStreamToFile(file, saveFilePath)

	if err != nil {
		callback.CallWithError(errors.New("can't read file"))
		return
	}

	callback.CallWithUndefined()
}

func JsRequestReadFormFileAsync(resHttpRequest *progpAPI.SharedResource, fieldName string, fileOffset int, callback progpAPI.ScriptFunction) {
	file, ok := openRequestFormFile(resHttpRequest, fieldName, fileOffset, callback)
	if !ok {
		return
	}

	defer file.Close()

	var buffer bytes.Buffer

	_, err := buffer.ReadFrom(file)
	if err != nil {
		callback.CallWithError(errors.New("can't read file"))
		return
	}

	callback.CallWithArrayBuffer2(buffer.Bytes())
}

func JsRequestWildcards(resHttpRequest *progpAPI.SharedResource) (error, []string) {
	call, ok := resHttpRequest.Value.(libHttpServer.HttpRequest)
	if !ok {
		return errors.New("invalid resource"), nil
	}

	return nil, call.GetWildcards()
}

func JsRequestRemainingSegments(resHttpRequest *progpAPI.SharedResource) (error, []string) {
	call, ok := resHttpRequest.Value.(libHttpServer.HttpRequest)
	if !ok {
		return errors.New("invalid resource"), nil
	}

	return nil, call.GetRemainingSegment()
}

func JsRequestCookie(resHttpRequest *progpAPI.SharedResource, cookieName string) (error, map[string]any) {
	call, ok := resHttpRequest.Value.(libHttpServer.HttpRequest)
	if !ok {
		return errors.New("invalid resource"), nil
	}

	c, err := call.GetCookie(cookieName)
	if err != nil {
		return err, nil
	}
	if c == nil {
		return nil, nil
	}

	out := make(map[string]any)
	return nil, out
}

func JsRequestHeaders(resHttpRequest *progpAPI.SharedResource) (error, map[string]string) {
	call, ok := resHttpRequest.Value.(libHttpServer.HttpRequest)
	if !ok {
		return errors.New("invalid resource"), nil
	}

	return nil, call.GetHeaders()
}

func JsRequestSetHeader(resHttpRequest *progpAPI.SharedResource, key string, value string) error {
	call, ok := resHttpRequest.Value.(libHttpServer.HttpRequest)
	if !ok {
		return errors.New("invalid resource")
	}

	call.SetHeader(key, value)
	return nil
}

func JsRequestCookies(resHttpRequest *progpAPI.SharedResource) (error, map[string]map[string]any) {
	call, ok := resHttpRequest.Value.(libHttpServer.HttpRequest)
	if !ok {
		return errors.New("invalid resource"), nil
	}

	c, err := call.GetCookies()
	return err, c
}

func JsRequestSetCookie(resHttpRequest *progpAPI.SharedResource, key string, value string, options libHttpServer.HttpCookieOptions) error {
	call, ok := resHttpRequest.Value.(libHttpServer.HttpRequest)
	if !ok {
		return errors.New("invalid resource")
	}

	return call.SetCookie(key, value, options)
}
