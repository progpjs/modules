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
	"encoding/json"
	"errors"
	"github.com/progpjs/httpServer/v2"
	"github.com/progpjs/httpServer/v2/libFastHttpImpl"
	"github.com/progpjs/progpAPI/v2"
	"github.com/progpjs/progpjs/v2"
	"mime/multipart"
	"net/textproto"
	"os"
	"path"
	"strings"
	"sync"
)

var NoResponseSendError = errors.New("no response send")

type httpFormFile struct {
	Header textproto.MIMEHeader `json:"header"`
	Size   int64                `json:"size"`
	Id     int                  `json:"id"`
}

func registerExportedFunctions() {
	rg := progpjs.GetFunctionRegistry()
	myMod := rg.UseGoNamespace("github.com/progpjs/modules/v2/modHttp")
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

	group.AddFunction("requestCookies", "JsRequestCookies", JsRequestCookies)
	group.AddFunction("requestCookie", "JsRequestCookie", JsRequestCookie)
	group.AddFunction("requestHeaders", "JsRequestHeaders", JsRequestHeaders)
	group.AddFunction("responseSetHeader", "JsRequestSetHeader", JsRequestSetHeader)
	group.AddFunction("responseSetCookie", "JsRequestSetCookie", JsRequestSetCookie)

	group.AddFunction("sendFileAsIs", "JsSendFileAsIs", JsSendFileAsIs)
	group.AddFunction("sendFile", "JsSendFile", JsSendFile)
	group.AddFunction("proxyTo", "JsProxyTo", JsProxyTo)
	group.AddFunction("serveFiles", "JsServerFiles", JsServerFiles)

	group.AddAsyncFunction("gzipCompressFile", "JsGzipCompressFileAsync", JsGzipCompressFileAsync)
	group.AddAsyncFunction("brotliCompressFile", "JsBrotliCompressFileAsync", JsBrotliCompressFileAsync)
	group.AddAsyncFunction("fetch", "JsFetchAsync", JsFetchAsync)
}

// JsConfigureServer configure a server designed by his port.
// It does nothing if the server is already started, but returns false if the configuration can't be applied.
func JsConfigureServer(serverPort int, params httpServer.StartParams) bool {
	server := libFastHttpImpl.GetFastHttpServer(serverPort)

	if server.IsStarted() {
		return false
	}

	server.SetStartServerParams(params)
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
func JsVerbWithFunction(rc *progpAPI.SharedResourceContainer, resHost *progpAPI.SharedResource, verb string, requestPath string, callback progpAPI.JsFunction) error {
	host, ok := resHost.Value.(*httpServer.HttpHost)
	if !ok {
		return errors.New("invalid resource")
	}

	// Allows calling this function more than one time.
	callback.KeepAlive()

	host.VERB(verb, requestPath, func(call httpServer.HttpRequest) error {
		res := rc.NewSharedResource(call, nil)

		// Allows disposing before the host script exit.
		defer res.Dispose()

		callback.CallWithResource2(res)

		// Will block the call until a response is sent.
		call.WaitResponse()

		if !call.IsBodySend() {
			return NoResponseSendError
		}

		return nil
	})

	return nil
}

// JsReturnString set the response to returns.
func JsReturnString(resHttpRequest *progpAPI.SharedResource, responseCode int, contentType string, responseText string) error {
	call, ok := resHttpRequest.Value.(httpServer.HttpRequest)
	if !ok {
		return errors.New("invalid resource")
	}

	call.SetContentType(contentType)

	// This response will unlock the caller mutex and the response will be sent.
	call.ReturnString(responseCode, responseText)
	return nil
}

func JsRequestURI(resHttpRequest *progpAPI.SharedResource) (error, string) {
	call, ok := resHttpRequest.Value.(httpServer.HttpRequest)
	if !ok {
		return errors.New("invalid resource"), ""
	}

	return nil, call.URI()
}

func JsRequestPath(resHttpRequest *progpAPI.SharedResource) (error, string) {
	call, ok := resHttpRequest.Value.(httpServer.HttpRequest)
	if !ok {
		return errors.New("invalid resource"), ""
	}

	return nil, call.Path()
}

func JsRequestIP(resHttpRequest *progpAPI.SharedResource) (error, string) {
	call, ok := resHttpRequest.Value.(httpServer.HttpRequest)
	if !ok {
		return errors.New("invalid resource"), ""
	}

	return nil, call.RemoteIP()
}

func JsRequestMethod(resHttpRequest *progpAPI.SharedResource) (error, string) {
	call, ok := resHttpRequest.Value.(httpServer.HttpRequest)
	if !ok {
		return errors.New("invalid resource"), ""
	}

	return nil, call.GetMethodName()
}

func JsRequestHost(resHttpRequest *progpAPI.SharedResource) (error, string) {
	call, ok := resHttpRequest.Value.(httpServer.HttpRequest)
	if !ok {
		return errors.New("invalid resource"), ""
	}

	return nil, call.GetHost().GetHostName()
}

func JsRequestQueryArgs(resHttpRequest *progpAPI.SharedResource) (error, map[string]string) {
	call, ok := resHttpRequest.Value.(httpServer.HttpRequest)
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
	call, ok := resHttpRequest.Value.(httpServer.HttpRequest)
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

func openRequestFormFile(resHttpRequest *progpAPI.SharedResource, fieldName string, fileOffset int, callback progpAPI.JsFunction) (multipart.File, bool) {
	call, ok := resHttpRequest.Value.(httpServer.HttpRequest)
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

func JsRequestSaveFormFileAsync(resHttpRequest *progpAPI.SharedResource, fieldName string, fileOffset int, saveFilePath string, callback progpAPI.JsFunction) {
	file, ok := openRequestFormFile(resHttpRequest, fieldName, fileOffset, callback)
	if !ok {
		return
	}

	defer func() {
		_ = file.Close()
	}()

	err := httpServer.SaveStreamBodyToFile(file, saveFilePath)

	if err != nil {
		callback.CallWithError(errors.New("can't read file"))
		return
	}

	callback.CallWithUndefined()
}

func JsRequestReadFormFileAsync(resHttpRequest *progpAPI.SharedResource, fieldName string, fileOffset int, callback progpAPI.JsFunction) {
	file, ok := openRequestFormFile(resHttpRequest, fieldName, fileOffset, callback)
	if !ok {
		return
	}

	defer func() {
		_ = file.Close()
	}()

	var buffer bytes.Buffer

	_, err := buffer.ReadFrom(file)
	if err != nil {
		callback.CallWithError(errors.New("can't read file"))
		return
	}

	callback.CallWithArrayBuffer2(buffer.Bytes())
}

func JsRequestWildcards(resHttpRequest *progpAPI.SharedResource) (error, []string) {
	call, ok := resHttpRequest.Value.(httpServer.HttpRequest)
	if !ok {
		return errors.New("invalid resource"), nil
	}

	return nil, call.GetWildcards()
}

func JsRequestCookie(resHttpRequest *progpAPI.SharedResource, cookieName string) (error, map[string]any) {
	call, ok := resHttpRequest.Value.(httpServer.HttpRequest)
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
	call, ok := resHttpRequest.Value.(httpServer.HttpRequest)
	if !ok {
		return errors.New("invalid resource"), nil
	}

	return nil, call.GetHeaders()
}

func JsRequestSetHeader(resHttpRequest *progpAPI.SharedResource, key string, value string) error {
	call, ok := resHttpRequest.Value.(httpServer.HttpRequest)
	if !ok {
		return errors.New("invalid resource")
	}

	call.SetHeader(key, value)
	return nil
}

func JsRequestCookies(resHttpRequest *progpAPI.SharedResource) (error, map[string]map[string]any) {
	call, ok := resHttpRequest.Value.(httpServer.HttpRequest)
	if !ok {
		return errors.New("invalid resource"), nil
	}

	c, err := call.GetCookies()
	return err, c
}

func JsRequestSetCookie(resHttpRequest *progpAPI.SharedResource, key string, value string, options httpServer.HttpCookieOptions) error {
	call, ok := resHttpRequest.Value.(httpServer.HttpRequest)
	if !ok {
		return errors.New("invalid resource")
	}

	return call.SetCookie(key, value, options)
}

func JsSendFileAsIs(resHttpRequest *progpAPI.SharedResource, filePath string, mimeType string, contentEncoding string) error {
	call, ok := resHttpRequest.Value.(httpServer.HttpRequest)
	if !ok {
		return errors.New("invalid resource")
	}

	return call.SendFileAsIs(filePath, mimeType, contentEncoding)
}

func JsSendFile(resHttpRequest *progpAPI.SharedResource, filePath string) error {
	call, ok := resHttpRequest.Value.(httpServer.HttpRequest)
	if !ok {
		return errors.New("invalid resource")
	}

	return call.SendFile(filePath)
}

func JsGzipCompressFileAsync(sourceFilePath string, destFilePath string, compressionLevel int, callback progpAPI.JsFunction) {
	progpAPI.SafeGoRoutine(func() {
		outDir := path.Dir(destFilePath)

		err := os.MkdirAll(outDir, os.ModePerm)
		if err != nil {
			callback.CallWithError(err)
			return
		}

		err = libFastHttpImpl.GzipCompressFile(sourceFilePath, destFilePath, compressionLevel)
		if err != nil {
			callback.CallWithError(err)
			return
		}

		callback.CallWithUndefined()
	})
}

func JsBrotliCompressFileAsync(sourceFilePath string, destFilePath string, compressionLevel int, callback progpAPI.JsFunction) {
	progpAPI.SafeGoRoutine(func() {
		outDir := path.Dir(destFilePath)

		err := os.MkdirAll(outDir, os.ModePerm)
		if err != nil {
			callback.CallWithError(err)
			return
		}

		err = libFastHttpImpl.BrotliCompressFile(sourceFilePath, destFilePath, compressionLevel)
		if err != nil {
			callback.CallWithError(err)
			return
		}

		callback.CallWithUndefined()
	})
}

func JsFetchAsync(url string, options JsFetchOptions, callback progpAPI.JsFunction) {
	progpAPI.SafeGoRoutine(func() {
		if options.Method == "" {
			options.Method = "GET"
		}

		fetchOptions := libFastHttpImpl.FetchOptions{
			SendHeaders: options.SendHeaders,
			SendCookies: options.SendCookies,
			SkipBody:    options.SkipBody,
			ContentType: options.ContentType,
			UserAgent:   options.UserAgent,
		}

		httpResult, err := libFastHttpImpl.Fetch(url, options.Method, fetchOptions)
		if err != nil {
			callback.CallWithError(err)
			return
		}

		defer httpResult.Dispose()

		jsResult := JsFetchResult{}
		jsResult.StatusCode = httpResult.StatusCode()

		if !options.SkipBody && ((jsResult.StatusCode == 200) || options.ForceReturningBody) {
			if options.StreamBodyToFile == "" {
				jsResult.Body, err = httpResult.GetBodyAsString()
				if err != nil {
					callback.CallWithError(err)
					return
				}
			} else {
				err = httpResult.StreamBodyToFile(options.StreamBodyToFile)
				if err != nil {
					callback.CallWithError(err)
					return
				}
			}
		}

		if options.ReturnHeaders {
			jsResult.Headers = httpResult.GetHeaders()
		}

		if options.ReturnCookies {
			jsResult.Cookies, err = httpResult.GetCookies()
			if err != nil {
				callback.CallWithError(err)
				return
			}
		}

		asJson, err := json.Marshal(jsResult)
		if err != nil {
			callback.CallWithError(err)
			return
		}

		callback.CallWithStringBuffer2(asJson)
	})
}

// JsProxyTo allows to proxy the incoming call directly to a website.
func JsProxyTo(resHost *progpAPI.SharedResource, requestPath string, targetHostName string, options JsProxyOptions) error {
	host, ok := resHost.Value.(*httpServer.HttpHost)
	if !ok {
		return errors.New("invalid resource")
	}

	mdw, err := libFastHttpImpl.BuildProxyAsIsMiddleware(targetHostName, 60)
	if err != nil {
		return err
	}

	host.AllVerbs(requestPath, mdw)

	if !options.ExcludeSubPaths {
		if !strings.HasPrefix(requestPath, "/") {
			requestPath += "/*"
		} else {
			requestPath += "*"
		}

		mdw, err := libFastHttpImpl.BuildProxyAsIsMiddleware(targetHostName, 60)
		if err != nil {
			return err
		}

		host.AllVerbs(requestPath, mdw)
	}

	return nil
}

func JsServerFiles(resHost *progpAPI.SharedResource, requestPath string, dirPath string, options JsServeFilesOptions) error {
	host, ok := resHost.Value.(*httpServer.HttpHost)
	if !ok {
		return errors.New("invalid resource")
	}

	if requestPath == "" {
		requestPath = "/"
	}

	mdw, err := libFastHttpImpl.BuildStaticFileServerMiddleware(requestPath, dirPath, libFastHttpImpl.StaticFileServerOptions{})
	if err != nil {
		return err
	}
	host.GET(requestPath, mdw)
	host.HEAD(requestPath, mdw)

	if requestPath[len(requestPath)-1] != '/' {
		requestPath += "/*"
	} else {
		requestPath += "*"
	}

	host.GET(requestPath, mdw)
	host.HEAD(requestPath, mdw)

	return nil
}

type JsFetchResult struct {
	StatusCode int                       `json:"statusCode"`
	Body       string                    `json:"body"`
	Headers    map[string]string         `json:"headers"`
	Cookies    map[string]map[string]any `json:"cookies"`
}

type JsFetchOptions struct {
	Method           string `json:"method"`
	StreamBodyToFile string `json:"streamBodyToFile"`
	ReturnHeaders    bool   `json:"returnHeaders"`
	ReturnCookies    bool   `json:"returnCookies"`

	SendHeaders map[string]string `json:"sendHeaders"`
	SendCookies map[string]string `json:"sendCookies"`

	// ForceReturningBody allows to return body event if response code isn't 200 Ok.
	ForceReturningBody bool `json:"forceReturningBody"`

	// SkipBody allows to avoid requesting the body.
	// Is useful if testing target existence or when we want to only get his headers.
	SkipBody bool `json:"skipBody"`

	// ContentType set the content type used when sending a body with the request.
	// Isn't set when no request are set.
	ContentType string `json:"contentType"`

	// UserAgent set the user agent used when sending a body with the request.
	// Isn't set when no request are set.
	UserAgent string
}

type JsProxyOptions struct {
	ExcludeSubPaths bool `json:"excludeSubPaths"`
}

type JsServeFilesOptions struct {
}
