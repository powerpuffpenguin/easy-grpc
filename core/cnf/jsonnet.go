package cnf

import (
	"encoding/json"

	"github.com/google/go-jsonnet"
)

func Load(filename string, v any) (e error) {
	vm := jsonnet.MakeVM()
	jsonStr, e := vm.EvaluateFile(filename)
	if e != nil {
		return
	}
	e = json.Unmarshal([]byte(jsonStr), v)
	if e != nil {
		return
	}
	return
}
