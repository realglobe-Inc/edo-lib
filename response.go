package util

import (
	"encoding/json"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"net/http"
	"strconv"
)

func JsonStringEscape(s string) string {
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

// {
//   "status": 500,
//   "message": "nani mo kamo oshimai dayo",
// }
func ErrorToResponseJson(err error) []byte {
	resp := errorToResponse(err)
	buff, err := json.Marshal(resp)
	if err != nil {
		err = erro.Wrap(err)
		log.Err("Json marshaling failed: ", resp)

		log.Err(erro.Unwrap(err))
		log.Debug(err)

		// 最後の手段。
		return []byte(`{"status":` + strconv.Itoa(http.StatusInternalServerError) + `,` +
			`"message":"` + JsonStringEscape(erro.Unwrap(err).Error()) + `"}`)
	}
	return buff
}

func errorToResponse(err error) interface{} {
	var resp struct {
		Stat int    `json:"status"`
		Msg  string `json:"message"`
	}

	switch e := erro.Unwrap(err).(type) {
	case *HttpStatusError:
		resp.Stat = e.Status()
		resp.Msg = e.Message()
	default:
		resp.Stat = http.StatusInternalServerError
		resp.Msg = e.Error()
	}
	return &resp
}
