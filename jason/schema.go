// Code generated by github.com/atombender/go-jsonschema, DO NOT EDIT.

package main

import "encoding/json"
import "fmt"
import "reflect"

type AnchorString string

type Applicator struct {
	// AdditionalProperties corresponds to the JSON schema field
	// "additionalProperties".
	AdditionalProperties interface{} `json:"additionalProperties,omitempty" yaml:"additionalProperties,omitempty" mapstructure:"additionalProperties,omitempty"`

	// AllOf corresponds to the JSON schema field "allOf".
	AllOf SchemaArray `json:"allOf,omitempty" yaml:"allOf,omitempty" mapstructure:"allOf,omitempty"`

	// AnyOf corresponds to the JSON schema field "anyOf".
	AnyOf SchemaArray `json:"anyOf,omitempty" yaml:"anyOf,omitempty" mapstructure:"anyOf,omitempty"`

	// Contains corresponds to the JSON schema field "contains".
	Contains interface{} `json:"contains,omitempty" yaml:"contains,omitempty" mapstructure:"contains,omitempty"`

	// DependentSchemas corresponds to the JSON schema field "dependentSchemas".
	DependentSchemas ApplicatorDependentSchemas `json:"dependentSchemas,omitempty" yaml:"dependentSchemas,omitempty" mapstructure:"dependentSchemas,omitempty"`

	// Else corresponds to the JSON schema field "else".
	Else interface{} `json:"else,omitempty" yaml:"else,omitempty" mapstructure:"else,omitempty"`

	// If corresponds to the JSON schema field "if".
	If interface{} `json:"if,omitempty" yaml:"if,omitempty" mapstructure:"if,omitempty"`

	// Items corresponds to the JSON schema field "items".
	Items *Schema `json:"items,omitempty" yaml:"items,omitempty" mapstructure:"items,omitempty"`

	// MaxContains corresponds to the JSON schema field "maxContains".
	MaxContains *NonNegativeInteger `json:"maxContains,omitempty" yaml:"maxContains,omitempty" mapstructure:"maxContains,omitempty"`

	// MinContains corresponds to the JSON schema field "minContains".
	MinContains NonNegativeInteger `json:"minContains,omitempty" yaml:"minContains,omitempty" mapstructure:"minContains,omitempty"`

	// Not corresponds to the JSON schema field "not".
	Not interface{} `json:"not,omitempty" yaml:"not,omitempty" mapstructure:"not,omitempty"`

	// OneOf corresponds to the JSON schema field "oneOf".
	OneOf SchemaArray `json:"oneOf,omitempty" yaml:"oneOf,omitempty" mapstructure:"oneOf,omitempty"`

	// PatternProperties corresponds to the JSON schema field "patternProperties".
	PatternProperties ApplicatorPatternProperties `json:"patternProperties,omitempty" yaml:"patternProperties,omitempty" mapstructure:"patternProperties,omitempty"`

	// PrefixItems corresponds to the JSON schema field "prefixItems".
	PrefixItems SchemaArray `json:"prefixItems,omitempty" yaml:"prefixItems,omitempty" mapstructure:"prefixItems,omitempty"`

	// Properties corresponds to the JSON schema field "properties".
	Properties ApplicatorProperties `json:"properties,omitempty" yaml:"properties,omitempty" mapstructure:"properties,omitempty"`

	// PropertyDependencies corresponds to the JSON schema field
	// "propertyDependencies".
	PropertyDependencies ApplicatorPropertyDependencies `json:"propertyDependencies,omitempty" yaml:"propertyDependencies,omitempty" mapstructure:"propertyDependencies,omitempty"`

	// PropertyNames corresponds to the JSON schema field "propertyNames".
	PropertyNames interface{} `json:"propertyNames,omitempty" yaml:"propertyNames,omitempty" mapstructure:"propertyNames,omitempty"`

	// Then corresponds to the JSON schema field "then".
	Then interface{} `json:"then,omitempty" yaml:"then,omitempty" mapstructure:"then,omitempty"`
}

type ApplicatorDependentSchemas map[string]Schema

type ApplicatorPatternProperties map[string]Schema

type ApplicatorProperties map[string]Schema

type ApplicatorPropertyDependencies map[string]map[string]Schema

// UnmarshalJSON implements json.Unmarshaler.
func (j *Applicator) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	type Plain Applicator
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	if v, ok := raw["dependentSchemas"]; !ok || v == nil {
		plain.DependentSchemas = ApplicatorDependentSchemas{}
	}
	if v, ok := raw["minContains"]; !ok || v == nil {
		plain.MinContains = 1.0
	}
	if v, ok := raw["patternProperties"]; !ok || v == nil {
		plain.PatternProperties = ApplicatorPatternProperties{}
	}
	if v, ok := raw["properties"]; !ok || v == nil {
		plain.Properties = ApplicatorProperties{}
	}
	if v, ok := raw["propertyDependencies"]; !ok || v == nil {
		plain.PropertyDependencies = ApplicatorPropertyDependencies{}
	}
	*j = Applicator(plain)
	return nil
}

type Content struct {
	// ContentEncoding corresponds to the JSON schema field "contentEncoding".
	ContentEncoding *string `json:"contentEncoding,omitempty" yaml:"contentEncoding,omitempty" mapstructure:"contentEncoding,omitempty"`

	// ContentMediaType corresponds to the JSON schema field "contentMediaType".
	ContentMediaType *string `json:"contentMediaType,omitempty" yaml:"contentMediaType,omitempty" mapstructure:"contentMediaType,omitempty"`

	// ContentSchema corresponds to the JSON schema field "contentSchema".
	ContentSchema interface{} `json:"contentSchema,omitempty" yaml:"contentSchema,omitempty" mapstructure:"contentSchema,omitempty"`
}

type Core struct {
	// Anchor corresponds to the JSON schema field "$anchor".
	Anchor *AnchorString `json:"$anchor,omitempty" yaml:"$anchor,omitempty" mapstructure:"$anchor,omitempty"`

	// Comment corresponds to the JSON schema field "$comment".
	Comment *string `json:"$comment,omitempty" yaml:"$comment,omitempty" mapstructure:"$comment,omitempty"`

	// Defs corresponds to the JSON schema field "$defs".
	Defs CoreDefs `json:"$defs,omitempty" yaml:"$defs,omitempty" mapstructure:"$defs,omitempty"`

	// DynamicAnchor corresponds to the JSON schema field "$dynamicAnchor".
	DynamicAnchor *AnchorString `json:"$dynamicAnchor,omitempty" yaml:"$dynamicAnchor,omitempty" mapstructure:"$dynamicAnchor,omitempty"`

	// DynamicRef corresponds to the JSON schema field "$dynamicRef".
	DynamicRef *IriReferenceString `json:"$dynamicRef,omitempty" yaml:"$dynamicRef,omitempty" mapstructure:"$dynamicRef,omitempty"`

	// Id corresponds to the JSON schema field "$id".
	Id *IriReferenceString `json:"$id,omitempty" yaml:"$id,omitempty" mapstructure:"$id,omitempty"`

	// Ref corresponds to the JSON schema field "$ref".
	Ref *IriReferenceString `json:"$ref,omitempty" yaml:"$ref,omitempty" mapstructure:"$ref,omitempty"`

	// Schema corresponds to the JSON schema field "$schema".
	Schema *IriString `json:"$schema,omitempty" yaml:"$schema,omitempty" mapstructure:"$schema,omitempty"`

	// Vocabulary corresponds to the JSON schema field "$vocabulary".
	Vocabulary CoreVocabulary `json:"$vocabulary,omitempty" yaml:"$vocabulary,omitempty" mapstructure:"$vocabulary,omitempty"`
}

type CoreDefs map[string]Schema

type CoreVocabulary map[string]bool

type FormatAnnotation struct {
	// Format corresponds to the JSON schema field "format".
	Format *string `json:"format,omitempty" yaml:"format,omitempty" mapstructure:"format,omitempty"`
}

type IriReferenceString string

type IriString string

type MetaData struct {
	// Default corresponds to the JSON schema field "default".
	Default interface{} `json:"default,omitempty" yaml:"default,omitempty" mapstructure:"default,omitempty"`

	// Deprecated corresponds to the JSON schema field "deprecated".
	Deprecated bool `json:"deprecated,omitempty" yaml:"deprecated,omitempty" mapstructure:"deprecated,omitempty"`

	// Description corresponds to the JSON schema field "description".
	Description *string `json:"description,omitempty" yaml:"description,omitempty" mapstructure:"description,omitempty"`

	// Examples corresponds to the JSON schema field "examples".
	Examples []interface{} `json:"examples,omitempty" yaml:"examples,omitempty" mapstructure:"examples,omitempty"`

	// ReadOnly corresponds to the JSON schema field "readOnly".
	ReadOnly bool `json:"readOnly,omitempty" yaml:"readOnly,omitempty" mapstructure:"readOnly,omitempty"`

	// Title corresponds to the JSON schema field "title".
	Title *string `json:"title,omitempty" yaml:"title,omitempty" mapstructure:"title,omitempty"`

	// WriteOnly corresponds to the JSON schema field "writeOnly".
	WriteOnly bool `json:"writeOnly,omitempty" yaml:"writeOnly,omitempty" mapstructure:"writeOnly,omitempty"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *MetaData) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	type Plain MetaData
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	if v, ok := raw["deprecated"]; !ok || v == nil {
		plain.Deprecated = false
	}
	if v, ok := raw["readOnly"]; !ok || v == nil {
		plain.ReadOnly = false
	}
	if v, ok := raw["writeOnly"]; !ok || v == nil {
		plain.WriteOnly = false
	}
	*j = MetaData(plain)
	return nil
}

type NonNegativeInteger int

type NonNegativeInteger_1 int

type Schema struct {
	Core `json:",inline"`

	Applicator `json:",inline"`

	Unevaluated `json:",inline"`

	Validation `json:",inline"`

	MetaData `json:",inline"`

	FormatAnnotation `json:",inline"`

	Content `json:",inline"`

	// RecursiveAnchor corresponds to the JSON schema field "$recursiveAnchor".
	RecursiveAnchor *bool `json:"$recursiveAnchor,omitempty" yaml:"$recursiveAnchor,omitempty" mapstructure:"$recursiveAnchor,omitempty"`

	// RecursiveRef corresponds to the JSON schema field "$recursiveRef".
	RecursiveRef *string `json:"$recursiveRef,omitempty" yaml:"$recursiveRef,omitempty" mapstructure:"$recursiveRef,omitempty"`

	// Definitions corresponds to the JSON schema field "definitions".
	Definitions SchemaDefinitions `json:"definitions,omitempty" yaml:"definitions,omitempty" mapstructure:"definitions,omitempty"`

	// Dependencies corresponds to the JSON schema field "dependencies".
	Dependencies SchemaDependencies `json:"dependencies,omitempty" yaml:"dependencies,omitempty" mapstructure:"dependencies,omitempty"`
}

type SchemaArray []Schema

type SchemaDefinitions map[string]Schema

type SchemaDependencies map[string]interface{}

// UnmarshalJSON implements json.Unmarshaler.
func (j *Schema) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	type Plain Schema
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	if v, ok := raw["definitions"]; !ok || v == nil {
		plain.Definitions = SchemaDefinitions{}
	}
	if v, ok := raw["dependencies"]; !ok || v == nil {
		plain.Dependencies = SchemaDependencies{}
	}
	*j = Schema(plain)
	return nil
}

type SimpleTypes string

const SimpleTypesArray SimpleTypes = "array"
const SimpleTypesBoolean SimpleTypes = "boolean"
const SimpleTypesInteger SimpleTypes = "integer"
const SimpleTypesNull SimpleTypes = "null"
const SimpleTypesNumber SimpleTypes = "number"
const SimpleTypesObject SimpleTypes = "object"
const SimpleTypesString SimpleTypes = "string"

var enumValues_SimpleTypes = []interface{}{
	"array",
	"boolean",
	"integer",
	"null",
	"number",
	"object",
	"string",
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *SimpleTypes) UnmarshalJSON(b []byte) error {
	var v string
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	var ok bool
	for _, expected := range enumValues_SimpleTypes {
		if reflect.DeepEqual(v, expected) {
			ok = true
			break
		}
	}
	if !ok {
		return fmt.Errorf("invalid value (expected one of %#v): %#v", enumValues_SimpleTypes, v)
	}
	*j = SimpleTypes(v)
	return nil
}

type StringArray []string

type Unevaluated struct {
	// UnevaluatedItems corresponds to the JSON schema field "unevaluatedItems".
	UnevaluatedItems interface{} `json:"unevaluatedItems,omitempty" yaml:"unevaluatedItems,omitempty" mapstructure:"unevaluatedItems,omitempty"`

	// UnevaluatedProperties corresponds to the JSON schema field
	// "unevaluatedProperties".
	UnevaluatedProperties interface{} `json:"unevaluatedProperties,omitempty" yaml:"unevaluatedProperties,omitempty" mapstructure:"unevaluatedProperties,omitempty"`
}

type Validation struct {
	// Const corresponds to the JSON schema field "const".
	Const interface{} `json:"const,omitempty" yaml:"const,omitempty" mapstructure:"const,omitempty"`

	// DependentRequired corresponds to the JSON schema field "dependentRequired".
	DependentRequired ValidationDependentRequired `json:"dependentRequired,omitempty" yaml:"dependentRequired,omitempty" mapstructure:"dependentRequired,omitempty"`

	// Enum corresponds to the JSON schema field "enum".
	Enum []interface{} `json:"enum,omitempty" yaml:"enum,omitempty" mapstructure:"enum,omitempty"`

	// ExclusiveMaximum corresponds to the JSON schema field "exclusiveMaximum".
	ExclusiveMaximum *float64 `json:"exclusiveMaximum,omitempty" yaml:"exclusiveMaximum,omitempty" mapstructure:"exclusiveMaximum,omitempty"`

	// ExclusiveMinimum corresponds to the JSON schema field "exclusiveMinimum".
	ExclusiveMinimum *float64 `json:"exclusiveMinimum,omitempty" yaml:"exclusiveMinimum,omitempty" mapstructure:"exclusiveMinimum,omitempty"`

	// MaxItems corresponds to the JSON schema field "maxItems".
	MaxItems *NonNegativeInteger_1 `json:"maxItems,omitempty" yaml:"maxItems,omitempty" mapstructure:"maxItems,omitempty"`

	// MaxLength corresponds to the JSON schema field "maxLength".
	MaxLength *NonNegativeInteger_1 `json:"maxLength,omitempty" yaml:"maxLength,omitempty" mapstructure:"maxLength,omitempty"`

	// MaxProperties corresponds to the JSON schema field "maxProperties".
	MaxProperties *NonNegativeInteger_1 `json:"maxProperties,omitempty" yaml:"maxProperties,omitempty" mapstructure:"maxProperties,omitempty"`

	// Maximum corresponds to the JSON schema field "maximum".
	Maximum *float64 `json:"maximum,omitempty" yaml:"maximum,omitempty" mapstructure:"maximum,omitempty"`

	// MinItems corresponds to the JSON schema field "minItems".
	MinItems ValidationMinItems `json:"minItems,omitempty" yaml:"minItems,omitempty" mapstructure:"minItems,omitempty"`

	// MinLength corresponds to the JSON schema field "minLength".
	MinLength ValidationMinLength `json:"minLength,omitempty" yaml:"minLength,omitempty" mapstructure:"minLength,omitempty"`

	// MinProperties corresponds to the JSON schema field "minProperties".
	MinProperties ValidationMinProperties `json:"minProperties,omitempty" yaml:"minProperties,omitempty" mapstructure:"minProperties,omitempty"`

	// Minimum corresponds to the JSON schema field "minimum".
	Minimum *float64 `json:"minimum,omitempty" yaml:"minimum,omitempty" mapstructure:"minimum,omitempty"`

	// MultipleOf corresponds to the JSON schema field "multipleOf".
	MultipleOf *float64 `json:"multipleOf,omitempty" yaml:"multipleOf,omitempty" mapstructure:"multipleOf,omitempty"`

	// Pattern corresponds to the JSON schema field "pattern".
	Pattern *string `json:"pattern,omitempty" yaml:"pattern,omitempty" mapstructure:"pattern,omitempty"`

	// Required corresponds to the JSON schema field "required".
	Required StringArray `json:"required,omitempty" yaml:"required,omitempty" mapstructure:"required,omitempty"`

	// Type corresponds to the JSON schema field "type".
	Type interface{} `json:"type,omitempty" yaml:"type,omitempty" mapstructure:"type,omitempty"`

	// UniqueItems corresponds to the JSON schema field "uniqueItems".
	UniqueItems bool `json:"uniqueItems,omitempty" yaml:"uniqueItems,omitempty" mapstructure:"uniqueItems,omitempty"`
}

type ValidationDependentRequired map[string]StringArray

type ValidationMinItems interface{}

type ValidationMinLength interface{}

type ValidationMinProperties interface{}

// UnmarshalJSON implements json.Unmarshaler.
func (j *Validation) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	type Plain Validation
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	if v, ok := raw["uniqueItems"]; !ok || v == nil {
		plain.UniqueItems = false
	}
	*j = Validation(plain)
	return nil
}
