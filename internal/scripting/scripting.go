package scripting

import (
	"context"
	"errors"
	"github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/stdlib"
	"ytb-downloader/internal/scripting/module"
)

type HandleResult struct {
	FilePath string
	Value    string
	Action   string
}

var moduleMap *tengo.ModuleMap

func init() {
	moduleMap = tengo.NewModuleMap()
	for _, name := range stdlib.AllModuleNames() {
		if mod := stdlib.BuiltinModules[name]; mod != nil {
			moduleMap.AddBuiltinModule(name, mod)
		}
		if mod := stdlib.SourceModules[name]; mod != "" {
			moduleMap.AddSourceModule(name, []byte(mod))
		}
	}
	moduleMap.Add("url", &tengo.BuiltinModule{
		Attrs: map[string]tengo.Object{
			"extractDomain": &tengo.UserFunction{
				Value: module.ExtractDomain,
			},
			"extractPath": &tengo.UserFunction{
				Value: module.ExtractPath,
			},
			"extractQuery": &tengo.UserFunction{
				Value: module.ExtractQuery,
			},
		},
	})
}

func HandleDownload(code []byte, input string) (HandleResult, error) {
	script := tengo.NewScript(code)
	_ = script.Add("_input", input)
	script.SetImports(moduleMap)

	compiled, err := script.RunContext(context.Background())
	if err != nil {
		return HandleResult{}, err
	}

	action := compiled.Get("_action")
	value := compiled.Get("_value")

	if action.IsUndefined() || value.IsUndefined() {
		return HandleResult{}, errors.New("_action or _value is undefined")
	}

	return HandleResult{
		Value:    value.String(),
		Action:   action.String(),
		FilePath: compiled.Get("_filepath").String(),
	}, nil
}
