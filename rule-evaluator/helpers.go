package rule_evaluator

import (
	"fmt"
	"reflect"
)

func (r *RuleEngine) CustomFn(Fn func(data interface{}) bool, description string) *SimpleCondition {
	return &SimpleCondition{
		Fn:          Fn,
		Description: description,
	}
}

// helpers to help you create simple conditions without defining functions every time, just give condition and description
// In these, you cannot use data argument passed to the function. You have to use closure in your caller function

func (r *RuleEngine) Not(condition bool, description string) *SimpleCondition {
	return r.CreateSimpleCondition(condition, description, true)
}

func (r *RuleEngine) Is(condition bool, description string) *SimpleCondition {
	return r.CreateSimpleCondition(condition, description, false)
}

// helpers for AND and OR conditions

func (r *RuleEngine) OR(conditions ...Condition) *ORCondition {
	return &ORCondition{
		Conditions: conditions,
	}
}

func (r *RuleEngine) AND(conditions ...Condition) *ANDCondition {
	return &ANDCondition{
		Conditions: conditions,
	}
}

// the following functions create the condition and description for you, given the field name in struct.
// but it can lead to errors if field value is not of expected type, which will only be known at runtime.

type ComparisonOperator string

const (
	GreaterThanOp          ComparisonOperator = ">"
	LessThanOp             ComparisonOperator = "<"
	GreaterThanOrEqualToOp ComparisonOperator = ">="
	LessThanOrEqualToOp    ComparisonOperator = "<="
	EqualToOp              ComparisonOperator = "=="
	NotEqualToOp           ComparisonOperator = "!="
	InOp                   ComparisonOperator = "in"
	IsNilOp                ComparisonOperator = "== nil"
	NotNilOp               ComparisonOperator = "!= nil"
)

func compareOpNumber(val interface{}, operand float64, op ComparisonOperator) bool {
	switch val.(type) {
	case int, int32, int64:
		intVal := int64(val.(int)) // Convert for consistency
		switch op {
		case GreaterThanOp:
			return intVal > int64(operand)
		case LessThanOp:
			return intVal < int64(operand)
		case GreaterThanOrEqualToOp:
			return intVal >= int64(operand)
		case LessThanOrEqualToOp:
			return intVal <= int64(operand)
		case EqualToOp:
			return intVal == int64(operand)
		case NotEqualToOp:
			return intVal != int64(operand)
		}
	case float32, float64:
		floatVal := val.(float64)
		switch op {
		case GreaterThanOp:
			return floatVal > operand
		case LessThanOp:
			return floatVal < operand
		case GreaterThanOrEqualToOp:
			return floatVal >= operand
		case LessThanOrEqualToOp:
			return floatVal <= operand
		case EqualToOp:
			return floatVal == operand
		case NotEqualToOp:
			return floatVal != operand

		}
	default:
		panic(fmt.Sprintf("ERROR: unsupported type for comparison: %T", val)) // Or handle with errors
	}
	return false // Should not reach here
}

func sliceContainsStr(list []string, operand string) bool {
	for _, item := range list {
		if item == operand {
			return true
		}
	}
	return false
}

func compareOpStr(val interface{}, operand string, op ComparisonOperator) bool {
	switch val.(type) {
	case string:
		strVal := val.(string)
		switch op {
		case EqualToOp:
			return strVal == operand
		case NotEqualToOp:
			return strVal != operand

		}
	case []string:
		switch op {
		case InOp:
			sliceVal := val.([]string)
			return sliceContainsStr(sliceVal, operand)
		}
	default:
		panic(fmt.Sprintf("ERROR: unsupported type for comparison: %T", val))
	}
	return false
}

func (r *RuleEngine) compareNumber(field string, operand float64, op ComparisonOperator) *SimpleCondition {
	return r.CustomFn(func(data interface{}) bool {
		switch typedData := data.(type) {
		case map[string]interface{}:
			if val, ok := typedData[field]; ok {
				return compareOpNumber(val, operand, op)
			} else {
				panic(fmt.Sprintf("ERROR: field '%s' not found in map", field))
			}
		default:
			fieldValue := reflect.ValueOf(data).FieldByName(field)
			if fieldValue.IsValid() {
				return compareOpNumber(fieldValue.Interface(), operand, op)
			} else {
				panic(fmt.Sprintf("ERROR: field '%s' not found in struct", field))
			}
		}
	}, fmt.Sprintf("%s %s %.2f", field, op, operand))
}

func (r *RuleEngine) compareString(field string, operand string, op ComparisonOperator) *SimpleCondition {
	return r.CustomFn(func(data interface{}) bool {
		switch typedData := data.(type) {
		case map[string]interface{}:
			if val, ok := typedData[field]; ok {
				return compareOpStr(val, operand, op)
			} else {
				panic(fmt.Sprintf("ERROR: field '%s' not found in map", field))
			}
		default:
			fieldValue := reflect.ValueOf(data).FieldByName(field)
			if fieldValue.IsValid() {
				return compareOpStr(fieldValue.Interface(), operand, op)
			} else {
				panic(fmt.Sprintf("ERROR: field '%s' not found in struct", field))
			}
		}
	}, fmt.Sprintf("%s %s %s", field, op, operand))
}

func (r *RuleEngine) checkNil(field string, invert bool, op ComparisonOperator) *SimpleCondition {
	return r.CustomFn(func(data interface{}) bool {
		switch typedData := data.(type) {
		case map[string]interface{}:
			if val, ok := typedData[field]; ok {
				if invert {
					return val != nil
				}
				return val == nil
			} else {
				panic(fmt.Sprintf("ERROR: field '%s' not found in map", field))
			}
		default:
			fieldValue := reflect.ValueOf(data).FieldByName(field)
			if fieldValue.IsValid() {
				if invert {
					return !fieldValue.IsNil()
				}
				return fieldValue.IsNil()
			} else {
				panic(fmt.Sprintf("ERROR: field '%s' not found in map", field))
			}
		}
		return false
	}, fmt.Sprintf("%s %s", field, op))
}

func (r *RuleEngine) GreaterThan(field string, value float64) *SimpleCondition {
	return r.compareNumber(field, value, GreaterThanOp)
}

func (r *RuleEngine) GreaterThanEqualTo(field string, value float64) *SimpleCondition {
	return r.compareNumber(field, value, GreaterThanOrEqualToOp)
}

func (r *RuleEngine) LessThan(field string, value float64) *SimpleCondition {
	return r.compareNumber(field, value, LessThanOp)
}

func (r *RuleEngine) LessThanEqualTo(field string, value float64) *SimpleCondition {
	return r.compareNumber(field, value, LessThanOrEqualToOp)
}

func (r *RuleEngine) EqualTo(field string, value float64) *SimpleCondition {
	return r.compareNumber(field, value, EqualToOp)
}

func (r *RuleEngine) NotEqualTo(field string, value float64) *SimpleCondition {
	return r.compareNumber(field, value, NotEqualToOp)
}

func (r *RuleEngine) EqualToStr(field string, value string) *SimpleCondition {
	return r.compareString(field, value, EqualToOp)
}

func (r *RuleEngine) NotEqualToStr(field string, value string) *SimpleCondition {
	return r.compareString(field, value, NotEqualToOp)
}

func (r *RuleEngine) ListContainsStr(field string, value string) *SimpleCondition {
	return r.compareString(field, value, InOp)
}

func (r *RuleEngine) IsNil(field string) *SimpleCondition {
	return r.checkNil(field, false, IsNilOp)
}

func (r *RuleEngine) IsNotNil(field string) *SimpleCondition {
	return r.checkNil(field, true, NotNilOp)
}
