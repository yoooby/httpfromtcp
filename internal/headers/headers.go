package headers

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
)

type Headers map[string]string

func (h Headers) GET(key string) string {
    key = strings.ToLower(key)
    if val, ok := h[key]; ok {
        return val
    }
    return ""
}

func NewHeaders() Headers {
	return Headers{}
}

var SEPERATOR = []byte("\r\n")
var ERROR_MALFORMED_DATA = fmt.Errorf("Malformed Headers")

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	var validHeaderKey = regexp.MustCompile(`^[!#$%&'*+\-.^_` + "`" + `|~0-9A-Za-z]+$`)

	idx := bytes.Index(data, SEPERATOR)
	if idx == -1 {
		return 0, false, nil
	}
	// check if it's the end of the headers
	if idx == 0 {
		return len(SEPERATOR), true, nil
	}	
	line := data[:idx]
	// trim whitespaces
	kv := bytes.SplitN(line, []byte(":"), 2)
	if len(kv) != 2 {
		return 0, false, ERROR_MALFORMED_DATA
	}
	fieldName := kv[0]
	fieldValue := kv[1]
	if len(fieldName) == 0 || bytes.HasSuffix(fieldName, []byte(" ")) {
		return 0, false, ERROR_MALFORMED_DATA
	}
	if !validHeaderKey.Match(fieldName) {
		return 0, false, fmt.Errorf("Invalid Character in header key: %s", fieldName)
	}

	key := strings.ToLower(string(bytes.TrimSpace(fieldName)))
	fieldValue = bytes.TrimSpace(fieldValue)
	println(key)
	if exisiting, ok := h[key]; ok {
		h[key] = exisiting + ", " + string(fieldValue)
	} else {
		h[key] = string(fieldValue)
	}

	return idx + len(SEPERATOR), false, nil

}




