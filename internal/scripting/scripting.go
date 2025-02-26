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
	Title    string
	Url      string
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
	url := compiled.Get("_url")

	if action.IsUndefined() || url.IsUndefined() {
		return HandleResult{}, errors.New("_action or _url is undefined")
	}

	filepath := compiled.Get("_filepath")

	if action.String() == "custom" && filepath.IsUndefined() {
		return HandleResult{}, errors.New("_filepath is required in custom mode")
	}

	return HandleResult{
		Url:      url.String(),
		Action:   action.String(),
		FilePath: getOptional(compiled, "_filepath"),
		Title:    getOptional(compiled, "_title"),
	}, nil
}

func getOptional(compiled *tengo.Compiled, variable string) string {
	if compiled.Get(variable).IsUndefined() {
		return ""
	}
	return compiled.Get(variable).String()
}
