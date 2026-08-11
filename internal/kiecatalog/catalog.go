// Copyright 2026 Kelvin Cushman and contributors. Licensed under Apache-2.0. See LICENSE.

// Package kiecatalog exposes the Kie Market model contracts captured from the
// official documentation index. The registry is embedded so local CLI and MCP
// discovery never spends API credits and never requires network access.
package kiecatalog

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"sort"
	"strings"
	"sync"
)

//go:embed catalog.json
var catalogJSON []byte

type Registry struct {
	SchemaVersion     int     `json:"schema_version"`
	SourceIndex       string  `json:"source_index"`
	SourceIndexSHA256 string  `json:"source_index_sha256"`
	ModelCount        int     `json:"model_count"`
	Models            []Model `json:"models"`
}

type Model struct {
	ID                  string         `json:"id"`
	Name                string         `json:"name"`
	Category            string         `json:"category"`
	Description         string         `json:"description,omitempty"`
	Source              string         `json:"source"`
	SourcePages         []string       `json:"source_pages,omitempty"`
	ModelIDSource       string         `json:"model_id_source,omitempty"`
	RequestSchema       map[string]any `json:"request_schema"`
	InputSchema         map[string]any `json:"input_schema"`
	RequestExample      any            `json:"request_example,omitempty"`
	InputExample        any            `json:"input_example,omitempty"`
	RequiredInputFields []string       `json:"required_input_fields"`
}

type Summary struct {
	ID                  string   `json:"id"`
	Name                string   `json:"name"`
	Category            string   `json:"category"`
	InputFields         []string `json:"input_fields"`
	RequiredInputFields []string `json:"required_input_fields"`
	Source              string   `json:"source"`
}

type ValidationIssue struct {
	Path    string `json:"path"`
	Keyword string `json:"keyword"`
	Message string `json:"message"`
}

var (
	loadOnce sync.Once
	loaded   Registry
	loadErr  error
	byID     map[string]int
)

func Load() (*Registry, error) {
	loadOnce.Do(func() {
		loadErr = json.Unmarshal(catalogJSON, &loaded)
		if loadErr != nil {
			loadErr = fmt.Errorf("decode embedded Kie model catalog: %w", loadErr)
			return
		}
		if loaded.ModelCount != len(loaded.Models) {
			loadErr = fmt.Errorf("embedded Kie model catalog count is %d, but contains %d models", loaded.ModelCount, len(loaded.Models))
			return
		}
		byID = make(map[string]int, len(loaded.Models))
		for i := range loaded.Models {
			if _, exists := byID[loaded.Models[i].ID]; exists {
				loadErr = fmt.Errorf("embedded Kie model catalog contains duplicate model %q", loaded.Models[i].ID)
				return
			}
			byID[loaded.Models[i].ID] = i
		}
	})
	if loadErr != nil {
		return nil, loadErr
	}
	return &loaded, nil
}

func Get(id string) (*Model, error) {
	registry, err := Load()
	if err != nil {
		return nil, err
	}
	id = strings.TrimSpace(id)
	index, ok := byID[id]
	if !ok {
		return nil, fmt.Errorf("Kie Market model %q is not in the embedded catalog; run the API refresh or use --allow-unknown-model for a newly released upstream model", id)
	}
	return &registry.Models[index], nil
}

func List(query, category string) ([]Summary, error) {
	registry, err := Load()
	if err != nil {
		return nil, err
	}
	query = strings.ToLower(strings.TrimSpace(query))
	category = strings.ToLower(strings.TrimSpace(category))
	models := make([]Summary, 0, len(registry.Models))
	for i := range registry.Models {
		model := &registry.Models[i]
		if category != "" && !strings.Contains(strings.ToLower(model.Category), category) {
			continue
		}
		haystack := strings.ToLower(strings.Join([]string{model.ID, model.Name, model.Category, model.Description}, " "))
		if query != "" && !strings.Contains(haystack, query) {
			continue
		}
		fields := schemaPropertyNames(model.InputSchema)
		models = append(models, Summary{
			ID: model.ID, Name: model.Name, Category: model.Category, InputFields: fields,
			RequiredInputFields: append([]string(nil), model.RequiredInputFields...), Source: model.Source,
		})
	}
	return models, nil
}

func Categories() ([]string, error) {
	registry, err := Load()
	if err != nil {
		return nil, err
	}
	seen := map[string]struct{}{}
	for _, model := range registry.Models {
		seen[model.Category] = struct{}{}
	}
	categories := make([]string, 0, len(seen))
	for category := range seen {
		categories = append(categories, category)
	}
	sort.Strings(categories)
	return categories, nil
}

func Example(id string) (map[string]any, error) {
	model, err := Get(id)
	if err != nil {
		return nil, err
	}
	if value, ok := model.InputExample.(map[string]any); ok && len(value) > 0 {
		return cloneMap(value), nil
	}
	value, _ := exampleValue(model.InputSchema, true).(map[string]any)
	if value == nil {
		value = map[string]any{}
	}
	return value, nil
}

func Validate(id string, input map[string]any) ([]ValidationIssue, error) {
	model, err := Get(id)
	if err != nil {
		return nil, err
	}
	issues := validateValue("input", input, model.InputSchema)
	sort.SliceStable(issues, func(i, j int) bool {
		if issues[i].Path == issues[j].Path {
			return issues[i].Keyword < issues[j].Keyword
		}
		return issues[i].Path < issues[j].Path
	})
	return issues, nil
}

func schemaPropertyNames(schema map[string]any) []string {
	properties, _ := schema["properties"].(map[string]any)
	fields := make([]string, 0, len(properties))
	for name := range properties {
		fields = append(fields, name)
	}
	sort.Strings(fields)
	return fields
}

func cloneMap(input map[string]any) map[string]any {
	data, _ := json.Marshal(input)
	var output map[string]any
	_ = json.Unmarshal(data, &output)
	return output
}

func exampleValue(schema map[string]any, required bool) any {
	if value, ok := schema["default"]; ok {
		return value
	}
	if value, ok := schema["example"]; ok {
		return value
	}
	if values, ok := schema["examples"].([]any); ok && len(values) > 0 {
		return values[0]
	}
	if values, ok := schema["enum"].([]any); ok && len(values) > 0 {
		return values[0]
	}
	switch schemaType(schema) {
	case "object":
		value := map[string]any{}
		properties, _ := schema["properties"].(map[string]any)
		requiredFields := stringSet(schema["required"])
		for name, raw := range properties {
			child, ok := raw.(map[string]any)
			if !ok {
				continue
			}
			_, isRequired := requiredFields[name]
			if isRequired || hasExampleValue(child) {
				value[name] = exampleValue(child, isRequired)
			}
		}
		return value
	case "array":
		if required {
			if items, ok := schema["items"].(map[string]any); ok {
				return []any{exampleValue(items, true)}
			}
		}
		return []any{}
	case "boolean":
		return false
	case "integer", "number":
		return 0
	case "string":
		if required {
			return "replace-me"
		}
		return ""
	default:
		return nil
	}
}

func hasExampleValue(schema map[string]any) bool {
	for _, key := range []string{"default", "example", "examples", "enum"} {
		if _, ok := schema[key]; ok {
			return true
		}
	}
	return false
}

func validateValue(path string, value any, schema map[string]any) []ValidationIssue {
	var issues []ValidationIssue
	for _, keyword := range []string{"allOf"} {
		if branches, ok := schema[keyword].([]any); ok {
			for _, raw := range branches {
				if branch, ok := raw.(map[string]any); ok {
					issues = append(issues, validateValue(path, value, branch)...)
				}
			}
		}
	}
	for _, keyword := range []string{"oneOf", "anyOf"} {
		if branches, ok := schema[keyword].([]any); ok && len(branches) > 0 {
			valid := 0
			for _, raw := range branches {
				if branch, ok := raw.(map[string]any); ok && len(validateValue(path, value, branch)) == 0 {
					valid++
				}
			}
			if (keyword == "oneOf" && valid != 1) || (keyword == "anyOf" && valid == 0) {
				issues = append(issues, issue(path, keyword, fmt.Sprintf("must match %s documented schema", keyword)))
			}
		}
	}
	typ := schemaType(schema)
	if typ != "" && !matchesType(value, typ) {
		return append(issues, issue(path, "type", fmt.Sprintf("must be %s, got %T", typ, value)))
	}
	if enum, ok := schema["enum"].([]any); ok && len(enum) > 0 {
		matched := false
		for _, candidate := range enum {
			if reflect.DeepEqual(normalizeNumber(value), normalizeNumber(candidate)) {
				matched = true
				break
			}
		}
		if !matched {
			issues = append(issues, issue(path, "enum", fmt.Sprintf("must be one of %s", compactJSON(enum))))
		}
	}
	switch typed := value.(type) {
	case map[string]any:
		properties, _ := schema["properties"].(map[string]any)
		for name := range stringSet(schema["required"]) {
			if _, ok := typed[name]; !ok {
				issues = append(issues, issue(joinPath(path, name), "required", "is required"))
			}
		}
		for name, childValue := range typed {
			if raw, ok := properties[name]; ok {
				if child, ok := raw.(map[string]any); ok {
					issues = append(issues, validateValue(joinPath(path, name), childValue, child)...)
				}
			} else if allowed, ok := schema["additionalProperties"].(bool); ok && !allowed {
				issues = append(issues, issue(joinPath(path, name), "additionalProperties", "is not a documented setting"))
			}
		}
	case []any:
		checkLength(&issues, path, len(typed), schema, "Items")
		if items, ok := schema["items"].(map[string]any); ok {
			for i, item := range typed {
				issues = append(issues, validateValue(fmt.Sprintf("%s[%d]", path, i), item, items)...)
			}
		}
	case string:
		checkLength(&issues, path, len([]rune(typed)), schema, "Length")
	case float64:
		if minimum, ok := number(schema["minimum"]); ok && typed < minimum {
			issues = append(issues, issue(path, "minimum", fmt.Sprintf("must be at least %v", minimum)))
		}
		if maximum, ok := number(schema["maximum"]); ok && typed > maximum {
			issues = append(issues, issue(path, "maximum", fmt.Sprintf("must be at most %v", maximum)))
		}
	}
	return issues
}

func checkLength(issues *[]ValidationIssue, path string, length int, schema map[string]any, suffix string) {
	if minimum, ok := number(schema["min"+suffix]); ok && float64(length) < minimum {
		*issues = append(*issues, issue(path, "min"+suffix, fmt.Sprintf("must contain at least %v", minimum)))
	}
	if maximum, ok := number(schema["max"+suffix]); ok && float64(length) > maximum {
		*issues = append(*issues, issue(path, "max"+suffix, fmt.Sprintf("must contain at most %v", maximum)))
	}
}

func schemaType(schema map[string]any) string {
	switch value := schema["type"].(type) {
	case string:
		return value
	case []any:
		for _, candidate := range value {
			if text, ok := candidate.(string); ok && text != "null" {
				return text
			}
		}
	}
	if _, ok := schema["properties"]; ok {
		return "object"
	}
	return ""
}

func matchesType(value any, typ string) bool {
	switch typ {
	case "object":
		_, ok := value.(map[string]any)
		return ok
	case "array":
		_, ok := value.([]any)
		return ok
	case "string":
		_, ok := value.(string)
		return ok
	case "boolean":
		_, ok := value.(bool)
		return ok
	case "number":
		_, ok := number(value)
		return ok
	case "integer":
		n, ok := number(value)
		return ok && math.Trunc(n) == n
	case "null":
		return value == nil
	default:
		return true
	}
}

func number(value any) (float64, bool) {
	switch typed := value.(type) {
	case float64:
		return typed, true
	case float32:
		return float64(typed), true
	case int:
		return float64(typed), true
	case int64:
		return float64(typed), true
	case json.Number:
		result, err := typed.Float64()
		return result, err == nil
	default:
		return 0, false
	}
}

func normalizeNumber(value any) any {
	if n, ok := number(value); ok {
		return n
	}
	return value
}

func stringSet(value any) map[string]struct{} {
	result := map[string]struct{}{}
	if values, ok := value.([]any); ok {
		for _, raw := range values {
			if text, ok := raw.(string); ok {
				result[text] = struct{}{}
			}
		}
	}
	return result
}

func joinPath(parent, child string) string {
	if parent == "" {
		return child
	}
	return parent + "." + child
}

func issue(path, keyword, message string) ValidationIssue {
	return ValidationIssue{Path: path, Keyword: keyword, Message: message}
}

func compactJSON(value any) string {
	data, _ := json.Marshal(value)
	return string(data)
}
