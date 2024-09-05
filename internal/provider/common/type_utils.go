package common

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// TypeFromValue converts an input to its type
func TypeFromValue(v interface{}) attr.Type {
	switch v := reflect.ValueOf(v); v.Kind() {
	case reflect.String:
		return types.StringType
	case reflect.Int:
		return types.Int32Type
	case reflect.Bool:
		return types.BoolType
	case reflect.Int32:
		return types.Int32Type
	case reflect.Int64:
		return types.Int64Type
	default:
		return types.StringType
	}
}

// ValueFromValue converts an input to its value
func ValueFromValue(v interface{}) attr.Value {
	switch v := reflect.ValueOf(v); v.Kind() {
	case reflect.String:
		return types.StringValue(v.String())
	case reflect.Int:
		return types.Int32Value(int32(v.Int()))
	case reflect.Bool:
		return types.BoolValue(v.Bool())
	case reflect.Int32:
		return types.Int32Value(int32(v.Int()))
	case reflect.Int64:
		return types.Int64Value(v.Int())
	default:
		return types.StringValue(v.String())
	}
}

// StructToAttrValues converts a struct to a map of attribute types and values.
// This iterates over the fields of the struct and extracts the values and types
func StructToAttrValues(s interface{}) (map[string]attr.Type, map[string]attr.Value) {
	t := reflect.TypeOf(s)
	v := reflect.ValueOf(s)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}

	typesMap := make(map[string]attr.Type)
	valuesMap := make(map[string]attr.Value)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("tfsdk")
		if tag == "" {
			continue
		}

		value := v.Field(i).Interface()
		typesMap[tag] = TypeFromValue(value)
		valuesMap[tag] = ValueFromValue(value)
	}

	return typesMap, valuesMap
}

// CheckNonNull checks if the required fields are not null or unknown.
func CheckNonNull(fields map[string]attr.Value, diags diag.Diagnostics, metadataName, typeName string) {
	optionalKeys := []string{}

	for field, value := range fields {
		if strings.HasSuffix(field, "?") {
			// The field is optional and will be checked later
			optionalKeys = append(optionalKeys, field)
			continue
		}

		if value.IsNull() {
			diags.AddError(
				fmt.Sprintf("%s is null", field),
				fmt.Sprintf("%s must be set for the %s %s.", field, metadataName, typeName),
			)
		}

		if value.IsUnknown() {
			diags.AddError(
				fmt.Sprintf("%s is unknown", field),
				fmt.Sprintf("%s must be set for the %s %s.", field, metadataName, typeName),
			)
		}
	}

	allNull := false
	allUnknown := false

	for _, key := range optionalKeys {
		if fields[key].IsNull() {
			allNull = true
		}

		if fields[key].IsUnknown() {
			allUnknown = true
		}
	}

	if allNull {
		diags.AddError(
			"all optional fields are null",
			fmt.Sprintf("At least one field must be set for the %s %s from %s", metadataName, typeName, strings.Join(optionalKeys, ", ")),
		)
	}

	if allUnknown {
		diags.AddError(
			"all optional fields are unknown",
			fmt.Sprintf("At least one field must be set for the %s %s from %s", metadataName, typeName, strings.Join(optionalKeys, ", ")),
		)
	}
}

// CheckRequired checks if the required fields are set.
func CheckRequired(ctx context.Context, fields map[string]attr.Value, diags diag.Diagnostics, metadataName, typeName string) {
	optionalKeys := []string{}

	for field, value := range fields {
		if strings.HasSuffix(field, "?") {
			// The field is optional and will be checked later
			optionalKeys = append(optionalKeys, field)
			continue
		}

		switch value.Type(ctx) {
		case basetypes.StringType{}:
			if value.(basetypes.StringValue).ValueString() == "" {
				diags.AddError(
					fmt.Sprintf("%s is empty", field),
					fmt.Sprintf("%s must be set for the %s %s.", field, metadataName, typeName),
				)
			}
		default:
			diags.AddError(
				fmt.Sprintf("%s is empty %s", field, value.Type(ctx)),
				fmt.Sprintf("%s must be set for the %s %s.", field, metadataName, typeName),
			)
		}
	}

	allEmpty := false

	for _, key := range optionalKeys {
		switch fields[key].Type(ctx) {
		case basetypes.StringType{}:
			if fields[key].(basetypes.StringValue).ValueString() == "" {
				allEmpty = true
			}
		default:
			allEmpty = true
		}
	}

	if allEmpty {
		diags.AddError(
			"all optional fields are empty",
			fmt.Sprintf("At least one field must be set for the %s %s from %s", metadataName, typeName, strings.Join(optionalKeys, ", ")),
		)
	}
}

// ToTypeValue converts a string to a types.String.
func ToTypeValue[K comparable](s K) (attr.Type, attr.Value) {
	switch s := reflect.ValueOf(s); s.Kind() {
	case reflect.String:
		return types.StringType, types.StringValue(s.String())
	case reflect.Int:
		return types.Int32Type, types.Int32Value(int32(s.Int()))
	case reflect.Bool:
		return types.BoolType, types.BoolValue(s.Bool())
	case reflect.Int32:
		return types.Int32Type, types.Int32Value(int32(s.Int()))
	case reflect.Int64:
		return types.Int64Type, types.Int64Value(s.Int())
	default:
		return types.StringType, types.StringValue(fmt.Sprintf("%v", s))
	}
}

// ToAttrList converts a list of strings to a []attr.Value
func ToAttrList[K comparable](list []K) []attr.Value {
	values := make([]attr.Value, len(list))

	for i, v := range list {
		_, values[i] = ToTypeValue(v)
	}

	return values
}

// ToListType converts a list of strings to a types.List.
func ToListType[K comparable, V attr.Type](list []K) (types.List, diag.Diagnostics) {
	// 2. Convert the list to an attr list
	values := ToAttrList(list)

	var typesType V

	// 3. Convert the attr list to a types.List
	return types.ListValue(typesType, values)
}

// FromListType converts a types.List to a list of strings.
func FromListType(ctx context.Context, list types.List) ([]string, diag.Diagnostics) {
	elements := make([]types.String, 0, len(list.Elements()))

	diags := list.ElementsAs(ctx, &elements, false)

	return FromStringList(elements), diags
}

// FromStringList converts a []types.String to a list of strings.
func FromStringList(list []types.String) []string {
	strs := make([]string, len(list))
	for i, s := range list {
		strs[i] = s.ValueString()
	}
	return strs
}
