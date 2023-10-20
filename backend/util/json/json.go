package jsonutil

import (
	"fmt"
	"strings"

	"github.com/valyala/fastjson"
)

type ErrMissingField struct {
	path []string
}

func (e ErrMissingField) Error() string {
	return fmt.Sprintf("missing field %q", strings.Join(e.path, "."))
}

type ErrInvalidField struct {
	path []string
	err  error
}

func (e ErrInvalidField) Error() string {
	return fmt.Sprintf("invalid field %q: %v", strings.Join(e.path, "."), e.err)
}

func Parse(data []byte) (Value, error) {
	jsonVal, err := fastjson.Parse(string(data))
	return Value{value: jsonVal}, err
}

type Value struct {
	value *fastjson.Value
	keys  []string
}

func (v Value) Get(keys ...string) Value {
	return Value{
		value: v.value.Get(keys...),
		keys:  keys,
	}
}

func (v Value) Bool() (bool, error) {
	if v.value == nil {
		return false, ErrMissingField{path: v.keys}
	}
	ret, err := v.value.Bool()
	if err != nil {
		return false, ErrInvalidField{path: v.keys, err: err}
	}
	return ret, nil
}

func (v Value) Int() (int, error) {
	if v.value == nil {
		return 0, ErrMissingField{path: v.keys}
	}
	ret, err := v.value.Int()
	if err != nil {
		return 0, ErrInvalidField{path: v.keys, err: err}
	}
	return ret, nil
}

func (v Value) Float() (float64, error) {
	if v.value == nil {
		return 0., ErrMissingField{path: v.keys}
	}
	ret, err := v.value.Float64()
	if err != nil {
		return 0., ErrInvalidField{path: v.keys, err: err}
	}
	return ret, nil
}

func (v Value) String() (string, error) {
	if v.value == nil {
		return "", ErrMissingField{path: v.keys}
	}
	ret, err := v.value.StringBytes()
	if err != nil {
		return "", ErrInvalidField{path: v.keys, err: err}
	}
	return string(ret), nil
}
