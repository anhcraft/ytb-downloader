package module

import (
	"errors"
	"github.com/d5/tengo/v2"
	"net/url"
)

func ExtractDomain(args ...tengo.Object) (tengo.Object, error) {
	if len(args) != 1 {
		return nil, tengo.ErrWrongNumArguments
	}

	input, ok := tengo.ToString(args[0])
	if !ok {
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "input",
			Expected: "string",
			Found:    args[0].TypeName(),
		}
	}

	parsedURL, err := url.Parse(input)
	if err != nil {
		return nil, errors.New("invalid url")
	}

	return &tengo.String{Value: parsedURL.Host}, nil
}

func ExtractPath(args ...tengo.Object) (tengo.Object, error) {
	if len(args) != 1 {
		return nil, tengo.ErrWrongNumArguments
	}

	input, ok := tengo.ToString(args[0])
	if !ok {
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "input",
			Expected: "string",
			Found:    args[0].TypeName(),
		}
	}

	parsedURL, err := url.Parse(input)
	if err != nil {
		return nil, errors.New("invalid url")
	}

	return &tengo.String{Value: parsedURL.Path}, nil
}

func ExtractQuery(args ...tengo.Object) (tengo.Object, error) {
	if len(args) != 1 {
		return nil, tengo.ErrWrongNumArguments
	}

	input, ok := tengo.ToString(args[0])
	if !ok {
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "input",
			Expected: "string",
			Found:    args[0].TypeName(),
		}
	}

	parsedURL, err := url.Parse(input)
	if err != nil {
		return nil, errors.New("invalid url")
	}

	queryParams, err := url.ParseQuery(parsedURL.RawQuery)
	if err != nil {
		return nil, errors.New("cannot decode query")
	}

	result := make(map[string]tengo.Object)
	for key, values := range queryParams {

		tengoValues := make([]tengo.Object, len(values))
		for i, val := range values {
			tengoValues[i] = &tengo.String{Value: val}
		}
		result[key] = &tengo.Array{Value: tengoValues}
	}

	return &tengo.ImmutableMap{Value: result}, nil
}
