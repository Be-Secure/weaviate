//                           _       _
// __      _____  __ ___   ___  __ _| |_ ___
// \ \ /\ / / _ \/ _` \ \ / / |/ _` | __/ _ \
//  \ V  V /  __/ (_| |\ V /| | (_| | ||  __/
//   \_/\_/ \___|\__,_| \_/ |_|\__,_|\__\___|
//
//  Copyright © 2016 - 2023 Weaviate B.V. All rights reserved.
//
//  CONTACT: hello@weaviate.io
//

package generate

type Params struct {
	Task           string
	ResultLanguage string
	OnSet          string
	Properties     []string
}

func (n Params) GetTask() string {
	return n.Task
}
func (n Params) GetResultLanguage() string {
	return n.ResultLanguage
}
func (n Params) GetOnSet() string {
	return n.OnSet
}
func (n Params) GetProperties() []string {
	return n.Properties
}
