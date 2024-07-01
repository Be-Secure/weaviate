//                           _       _
// __      _____  __ ___   ___  __ _| |_ ___
// \ \ /\ / / _ \/ _` \ \ / / |/ _` | __/ _ \
//  \ V  V /  __/ (_| |\ V /| | (_| | ||  __/
//   \_/\_/ \___|\__,_| \_/ |_|\__,_|\__\___|
//
//  Copyright © 2016 - 2024 Weaviate B.V. All rights reserved.
//
//  CONTACT: hello@weaviate.io
//

package generate

import (
	"log"
	"regexp"
	"strings"

	"github.com/weaviate/weaviate/entities/models"
	"github.com/weaviate/weaviate/entities/schema"

	"github.com/tailor-inc/graphql/language/ast"
)

var compile, _ = regexp.Compile(`{([\w\s]*?)}`)

func (p *GenerateProvider) parseGenerateArguments(args []*ast.Argument, class *models.Class) *Params {
	out := &Params{}

	propertiesToExtract := make([]string, 0)

	for _, arg := range args {
		switch arg.Name.Value {
		case "singleResult":
			obj := arg.Value.(*ast.ObjectValue).Fields
			out.Prompt = &obj[0].Value.(*ast.StringValue).Value
			singlePropPrompts := ExtractPropsFromPrompt(out.Prompt)
			propertiesToExtract = append(propertiesToExtract, singlePropPrompts...)
		case "groupedResult":
			obj := arg.Value.(*ast.ObjectValue).Fields
			propertiesProvided := false
			for _, field := range obj {
				switch field.Name.Value {
				case "task":
					out.Task = &field.Value.(*ast.StringValue).Value
				case "properties":
					inp := field.Value.GetValue().([]ast.Value)
					out.Properties = make([]string, len(inp))

					for i, value := range inp {
						out.Properties[i] = value.(*ast.StringValue).Value
					}
					propertiesToExtract = append(propertiesToExtract, out.Properties...)
					propertiesProvided = true
				}
			}
			if !propertiesProvided {
				propertiesToExtract = append(propertiesToExtract, schema.GetPropertyNamesFromClass(class, false)...)
			}

		default:
			// ignore what we don't recognize
			log.Printf("Igonore not recognized value: %v", arg.Name.Value)
		}
	}

	out.PropertiesToExtract = propertiesToExtract

	return out
}

func ExtractPropsFromPrompt(prompt *string) []string {
	propertiesToExtract := make([]string, 0)
	all := compile.FindAll([]byte(*prompt), -1)
	for entry := range all {
		propName := string(all[entry])
		propName = strings.Trim(propName, "{")
		propName = strings.Trim(propName, "}")
		propertiesToExtract = append(propertiesToExtract, propName)
	}
	return propertiesToExtract
}
