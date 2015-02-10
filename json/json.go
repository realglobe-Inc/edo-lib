package json

import ()

func StringEscape(s string) string {
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
