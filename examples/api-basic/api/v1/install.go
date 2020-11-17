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

package v1

import (
	"github.com/caicloud/nirvana/definition"
)

var descriptors = []definition.RESTfulDescriptor{}

func register(ds ...definition.RESTfulDescriptor) {
	descriptors = append(descriptors, ds...)
}

func Descriptor() definition.RESTfulDescriptor {
	return definition.RESTfulDescriptor{
		Path:        "/api/v1",
		Description: "It contains all APIs in v1",
		Consumes:    []string{definition.MIMEAll},
		Produces:    []string{definition.MIMEJSON},
		Children:    descriptors,
	}
}
