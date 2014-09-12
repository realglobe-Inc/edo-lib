package util

import (
	"encoding/json"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"reflect"
)

func jsonEscape(s string) string {
	output := ""
	for _, r := range s {
		switch r {
		case '"':
			output += "\\\""
		case '\\':
			output += "\\\\"
		case '/':
			output += "\\/"
		case '\n':
			output += "\\n"
		case '\r':
			output += "\\r"
		case '\t':
			output += "\\t"
		case '\b':
			output += "\\b"
		case '\f':
			output += "\\f"
		default:
			output += string(r)
		}
	}
	return output
}

func toResponseJson(res interface{}) []byte {
	output, err := json.Marshal(res)
	if err != nil {
		err = erro.Wrap(err)
		log.Err("Marshal failed: ", res)
		log.Err(erro.Unwrap(err))
		log.Debug(err)

		trc := err.(*erro.Tracer)
		t := reflect.TypeOf(trc.Cause())
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		sysType := t.Name()
		if sysType == "" {
			sysType = "Unknown"
		}

		// 最後の手段。
		return []byte(`{` +
			`"name":"Error",` +
			`"message":"` + jsonEscape(trc.Cause().Error()) + `",` +
			`"sys_type":""` + sysType + `",` +
			`"sys_data":{}` +
			`"sys_stack":"` + jsonEscape(trc.Stack()) + `"` +
			`}`)
	}
	return output
}

// edo-interpreter の出力にならう。
// {
//   "name": "Error",
//   "message": "hoge is not exist",
//   "stack": "...",
//   "sys_type": "eventNotFound",
//   "sys_data": {
//     "event": "hoge"
//   },
//   "sys_stack": "..."
// }
func ErrorToResponseJson(err error) []byte {
	return toResponseJson(errorToResponse(err))
}

func errorToResponse(err error) interface{} {
	raw := erro.Unwrap(err)

	switch e := raw.(type) {
	case *HttpStatusError:
		return httpStatusErrorToResponse(e)
	default:
		var res struct {
			Name     string      `json:"name"`
			Message  string      `json:"message"`
			SysType  string      `json:"sys_type"`
			SysData  interface{} `json:"sys_data"`
			SysStack string      `json:"sys_stack,omitempty"`
		}

		res.Name = "Error"
		res.Message = raw.Error()

		// sysType.
		t := reflect.TypeOf(raw)
		for t.Kind() == reflect.Ptr {
			t = t.Elem()
		}

		sysType := t.Name()
		if sysType == "" {
			sysType = "Unknown"
		}
		res.SysType = sysType

		res.SysData = raw

		// sysStack.
		var sysStack string
		switch r := raw.(type) {
		case *PanicWrapper:
			sysStack += r.stack
		}

		trc, ok := err.(*erro.Tracer)
		if ok {
			if len(sysStack) < 0 {
				sysStack += "\n"
			}
			sysStack += trc.Stack()
		}
		res.SysStack = sysStack

		return &res
	}
}

func httpStatusErrorToResponse(err *HttpStatusError) interface{} {
	var res struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	}

	res.Status = err.Status()
	res.Message = err.Message()

	return &res
}
