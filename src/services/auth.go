package services

import (
	"errors"
	"reflect"
)

// function list stored
var FunctionMap = map[string]interface{}{
	"Login": Login,
}

// CheckFunction dynamically calls a function
func Auth(str string) (string, error) {
	if check, err := CheckFunction(str); err != nil {
		return "", err
	} else {
		return check, nil
	}
}

func CheckFunction(funcName string) (string, error) {
	if fn, exists := FunctionMap[funcName]; exists {

		if reflect.ValueOf(fn).Kind() == reflect.Func {

			if f, ok := fn.(func() string); ok {
				return f(), nil
			}

			return "", errors.New("bad Request")
		}

		return "", errors.New("bad Request")
	}

	return "", errors.New("bad request")
}

func Login() string {
	return "Foo called"
}
