package modHttp

import (
	"github.com/progpjs/progpAPI/v2"
	"github.com/progpjs/progpjs/v2"
)

//region (SharedResource, StringBuffer)

var gCallJsFunctionWith_SharedResource_StringBuffer CallJsFunctionWith_SharedResource_StringBuffer

type CallJsFunctionWith_SharedResource_StringBuffer interface {
	Call(js progpAPI.JsFunction, res *progpAPI.SharedResource, value progpAPI.StringBuffer)
}

type impl__CallJsFunctionWith_SharedResource_StringBuffer struct {
}

func (*impl__CallJsFunctionWith_SharedResource_StringBuffer) Call(js progpAPI.JsFunction, res *progpAPI.SharedResource, value progpAPI.StringBuffer) {
	js.DynamicFunctionCaller(res, value)
}

//endregion

func registerJsFunctionCallers() {
	gCallJsFunctionWith_SharedResource_StringBuffer = progpjs.GetFunctionCaller(&impl__CallJsFunctionWith_SharedResource_StringBuffer{}).(CallJsFunctionWith_SharedResource_StringBuffer)
}
