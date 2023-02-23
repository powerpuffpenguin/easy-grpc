package cnf

import (
	"encoding/json"

	"github.com/google/go-jsonnet"
	"github.com/powerpuffpenguin/easy-grpc/core"
)

func Load(filename string, v any) (e error) {
	vm := jsonnet.MakeVM()
	jsonStr, e := vm.EvaluateFile(filename)
	if e != nil {
		return
	}

	e = json.Unmarshal(core.StringToBytes(jsonStr), v)
	if e != nil {
		return
	}
	return
}
