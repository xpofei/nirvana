/*
Copyright 2017 Caicloud Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package definition

import (
	"context"
	"fmt"
	"reflect"
)

// MIME types
const (
	// MIMEAll indicates a accept type from http request.
	// It means client can receive any content.
	// Request content type in header "Content-Type" must not set to "*/*".
	// It only can exist in request header "Accept".
	// In most time, it locate at the last element of "Accept".
	// It's default value if client have not set "Accept" header.
	MIMEAll         = "*/*"
	MIMENone        = ""
	MIMEText        = "text/plain"
	MIMEHTML        = "text/html"
	MIMEJSON        = "application/json"
	MIMEXML         = "application/xml"
	MIMEOctetStream = "application/octet-stream"
	MIMEURLEncoded  = "application/x-www-form-urlencoded"
	MIMEFormData    = "multipart/form-data"
)

// DataErrorResults returns the most frequently-used results.
// Definition function should have two results. The first is
// any type for data, and the last is error.
func DataErrorResults(description string) []Result {
	return []Result{
		{
			Destination: Data,
			Description: description,
		},
		{
			Destination: Error,
		},
	}
}

// ParameterFor creates a simple parameter.
func ParameterFor(source Source, name string, description string, operators ...Operator) Parameter {
	return Parameter{
		Source:      source,
		Name:        name,
		Description: description,
		Operators:   operators,
	}
}

// PathParameterFor creates a path parameter
func PathParameterFor(name string, description string, operators ...Operator) Parameter {
	return ParameterFor(Path, name, description, operators...)
}

// QueryParameterFor creates a query parameter
func QueryParameterFor(name string, description string, operators ...Operator) Parameter {
	return ParameterFor(Query, name, description, operators...)
}

// HeaderParameterFor creates a header parameter
func HeaderParameterFor(name string, description string, operators ...Operator) Parameter {
	return ParameterFor(Header, name, description, operators...)
}

// FormParameterFor creates a form parameter
func FormParameterFor(name string, description string, operators ...Operator) Parameter {
	return ParameterFor(Form, name, description, operators...)
}

// FileParameterFor creates a file parameter
func FileParameterFor(name string, description string, operators ...Operator) Parameter {
	return ParameterFor(File, name, description, operators...)
}

// BodyParameterFor creates a body parameter
func BodyParameterFor(description string, operators ...Operator) Parameter {
	return ParameterFor(Body, "", description, operators...)
}

// PrefabParameterFor creates a prefab parameter
func PrefabParameterFor(name string, description string, operators ...Operator) Parameter {
	return ParameterFor(Prefab, name, description, operators...)
}

// AutoParameterFor creates an auto parameter
func AutoParameterFor(description string, operators ...Operator) Parameter {
	return ParameterFor(Auto, "", description, operators...)
}

// ResultFor creates a simple result.
func ResultFor(dest Destination, description string, operators ...Operator) Result {
	return Result{
		Destination: dest,
		Description: description,
		Operators:   operators,
	}
}

// MetaResultFor creates meta result.
func MetaResultFor(description string, operators ...Operator) Result {
	return ResultFor(Meta, description, operators...)
}

// DataResultFor creates data result.
func DataResultFor(description string, operators ...Operator) Result {
	return ResultFor(Data, description, operators...)
}

// ErrorResult creates error result.
func ErrorResult() Result {
	return ResultFor(Error, "")
}

// OperatorFunc creates operator by function.
// function must has signature:
//  func f(context.Context, string, AnyType) (AnyType, error)
// The second parameter is a string that is used to identify field.
// AnyType can be any type in go. But struct type and
// built-in data type is recommended.
func OperatorFunc(kind string, f interface{}) Operator {
	typ := reflect.TypeOf(f)
	if typ.Kind() != reflect.Func {
		panic(fmt.Sprintf("Parameter f in OperatorFunc must be a function, but got: %s", typ.Kind()))
	}
	if typ.NumIn() != 3 {
		panic(fmt.Sprintf("Function must have 3 parameters, but got: %d", typ.NumIn()))
	}
	if typ.In(0) != reflect.TypeOf((*context.Context)(nil)).Elem() {
		panic(fmt.Sprintf("The first parameter of function must be context.Context, but got: %v", typ.In(0)))
	}
	if typ.In(1) != reflect.TypeOf("") {
		panic(fmt.Sprintf("The second parameter of function must be string, but got: %v", typ.In(0)))
	}
	if typ.NumOut() != 2 {
		panic(fmt.Sprintf("Parameter f in OperatorFunc must have two results, but got: %d", typ.NumOut()))
	}
	if typ.Out(1).String() != "error" {
		panic(fmt.Sprintf("The last result of parameter f in OperatorFunc must be error, but got: %v", typ.Out(1)))
	}
	return &operatorRef{
		kind:  kind,
		in:    typ.In(2),
		out:   typ.Out(0),
		value: reflect.ValueOf(f),
	}
}

// NewOperator creates operator by function.
// function must has signature:
//  func f(context.Context, AnyType) (AnyType, error)
// AnyType can be any type in go. But struct type and
// built-in data type is recommended.
func NewOperator(kind string, in, out reflect.Type, f func(ctx context.Context, field string, object interface{}) (interface{}, error)) Operator {
	return &operator{
		kind: kind,
		in:   in,
		out:  out,
		f:    f,
	}
}

type operator struct {
	kind string
	in   reflect.Type
	out  reflect.Type
	f    func(ctx context.Context, field string, object interface{}) (interface{}, error)
}

// Kind indicates operator type.
func (o *operator) Kind() string {
	return o.kind
}

// In returns the type of the only object parameter of operator.
func (o *operator) In() reflect.Type {
	return o.in
}

// Out returns the type of the only object result of operator.
func (o *operator) Out() reflect.Type {
	return o.out
}

// Operate operates an object and return one.
func (o *operator) Operate(ctx context.Context, field string, object interface{}) (interface{}, error) {
	return o.f(ctx, field, object)
}

type operatorRef struct {
	kind  string
	in    reflect.Type
	out   reflect.Type
	value reflect.Value
}

// Kind indicates operator type.
func (o *operatorRef) Kind() string {
	return o.kind
}

// In returns the type of the only object parameter of operator.
func (o *operatorRef) In() reflect.Type {
	return o.in
}

// Out returns the type of the only object result of operator.
func (o *operatorRef) Out() reflect.Type {
	return o.out
}

// Operate operates an object and return one.
func (o *operatorRef) Operate(ctx context.Context, field string, object interface{}) (interface{}, error) {
	var objectValue reflect.Value
	if object == nil {
		objectValue = reflect.Zero(o.in)
	} else {
		objectValue = reflect.ValueOf(object)
	}

	results := o.value.Call([]reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(field), objectValue})
	v := results[1]
	if v.IsNil() {
		return results[0].Interface(), nil
	}
	return nil, results[1].Interface().(error)
}

// SimpleDescriptor creates a simple REST descriptor for handler.
// The descriptor consumes all content types and produces all accept types.
func SimpleDescriptor(method Method, path string, f interface{}) Descriptor {
	return Descriptor{
		Path: path,
		Definitions: []Definition{
			{
				Method:   method,
				Function: f,
				Consumes: []string{MIMEAll},
				Produces: []string{MIMEAll},
			},
		},
	}
}

// SimpleRPCDescriptor creates a simple RPC descriptor for handler.
// The descriptor consumes all content types and produces all accept types.
func SimpleRPCDescriptor(path string, f interface{}) RPCDescriptor {
	return RPCDescriptor{
		Path: path,
		Actions: []RPCAction{
			{
				Function: f,
				Consumes: []string{MIMEAll},
				Produces: []string{MIMEAll},
			},
		},
	}
}
