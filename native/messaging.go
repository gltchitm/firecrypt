package native

import (
	"C"
	"encoding/base64"
	"encoding/json"
	"strings"
)

var onMessageCallback func(string, []string) interface{}

//export onMessage
func onMessage(payload *C.char) *C.char {
	splitData := strings.Split(C.GoString(payload), ",")

	name, err := base64.StdEncoding.DecodeString(splitData[0])
	if err != nil {
		panic(err)
	}

	detail, err := base64.StdEncoding.DecodeString(splitData[1])
	if err != nil {
		panic(err)
	}

	response, err := json.Marshal(onMessageCallback(
		string(name),
		strings.Split(string(detail), ","),
	))
	if err != nil {
		panic(err)
	}

	return C.CString(base64.StdEncoding.EncodeToString(response))
}
