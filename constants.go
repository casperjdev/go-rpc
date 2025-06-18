package main

import (
	"fmt"
	"maps"
	"reflect"
	"strings"
)

// Constants
// You can add as many constants as you want here and set them to anything
func getConstants(data map[string]string) map[string]string {
	constants := map[string]string{}

	// External commands from your config
    for placeholder, scriptPath := range data {
        output, err := runExternalScript(scriptPath)
        if err != nil {
            fmt.Printf("Error running external constant %s: %v\n", placeholder, err)
            output = placeholder
        } else {
					constants[placeholder] = output
					fmt.Printf("Successfully loaded constant %s: %v\n", placeholder, output)
				}
    }

	return constants
}

func MergeConstants(static, dynamic map[string]string) map[string]string {
	merged := make(map[string]string)

	maps.Copy(merged, static)
	maps.Copy(merged, dynamic)
	
	return merged
}

func ProcessConstants(v any, constants map[string]string) {
	val := reflect.ValueOf(v)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return
	}

	for i := range val.NumField() {
		field := val.Field(i)

		if !field.CanSet() {
			continue
		}

		switch field.Kind() {
		case reflect.String:
			original := field.String()
			processed := replaceConstants(original, constants)
			field.SetString(processed)

		case reflect.Struct:
			ProcessConstants(field.Addr().Interface(), constants)

		case reflect.Ptr:
			if !field.IsNil() && field.Elem().Kind() == reflect.Struct {
				ProcessConstants(field.Interface(), constants)
			}

		case reflect.Slice:
			if field.Type().Elem().Kind() == reflect.Struct {
				for j := 0; j < field.Len(); j++ {
					ProcessConstants(field.Index(j).Addr().Interface(), constants)
				}
			}
		}
	}
}

func replaceConstants(input string, constants map[string]string) string {
	for placeholder, value := range constants {
		input = strings.ReplaceAll(input, placeholder, value)
	}
	return input
}
