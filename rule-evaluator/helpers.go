package rule_evaluator

import (
	"fmt"
	"reflect"
)

func (r *RuleEngine[T]) CustomFn(Fn func(data T) bool, description string) *SimpleCondition[T] {
	return &SimpleCondition[T]{
		Fn:          Fn,
		Description: description,
	}
}

// helpers to help you create simple conditions without defining functions every time, just give condition and description
// In these, you cannot use data argument passed to the function. You have to use closure in your caller function

func (r *RuleEngine[T]) Not(condition bool, description string) *SimpleCondition[T] {
	return r.CreateSimpleCondition(condition, description, true)
}

func (r *RuleEngine[T]) Is(condition bool, description string) *SimpleCondition[T] {
	return r.CreateSimpleCondition(condition, description, false)
}

// helpers for AND and OR conditions

func (r *RuleEngine[T]) OR(conditions ...Condition[T]) *ORCondition[T] {
	return &ORCondition[T]{
		Conditions: conditions,
	}
}

func (r *RuleEngine[T]) AND(conditions ...Condition[T]) *ANDCondition[T] {
	return &ANDCondition[T]{
		Conditions: conditions,
	}
}

// the following functions create the condition and description for you, given the field name in struct.
// but it can lead to errors if field value is not of expected type, which will only be known at runtime.

func (r *RuleEngine[T]) GreaterThan(field string, operand float64) *SimpleCondition[T] {
	return r.CustomFn(func(data T) bool {
		fieldValue := reflect.ValueOf(data).FieldByName(field)
		switch fieldValue.Kind() {
		case reflect.Int, reflect.Int32, reflect.Int64:
			return fieldValue.Int() > int64(operand)
		case reflect.Float32, reflect.Float64:
			return fieldValue.Float() > operand
		default:
			panic(fmt.Sprintf("ERROR GreaterThan: unsupported type of field value for key %s", field))
		}
	}, fmt.Sprintf("%s > %.2f", field, operand))
}

func (r *RuleEngine[T]) GreaterThanEqualTo(field string, operand float64) *SimpleCondition[T] {
	return r.CustomFn(func(data T) bool {
		fieldValue := reflect.ValueOf(data).FieldByName(field)
		switch fieldValue.Kind() {
		case reflect.Int, reflect.Int32, reflect.Int64:
			return fieldValue.Int() >= int64(operand)
		case reflect.Float32, reflect.Float64:
			return fieldValue.Float() >= operand
		default:
			panic(fmt.Sprintf("ERROR GreaterThanEqualTo: unsupported type of field value for key %s", field))
		}
	}, fmt.Sprintf("%s >= %.2f", field, operand))
}

func (r *RuleEngine[T]) LessThan(field string, operand float64) *SimpleCondition[T] {
	return r.CustomFn(func(data T) bool {
		fieldValue := reflect.ValueOf(data).FieldByName(field)
		switch fieldValue.Kind() {
		case reflect.Int, reflect.Int32, reflect.Int64:
			return fieldValue.Int() < int64(operand)
		case reflect.Float32, reflect.Float64:
			return fieldValue.Float() < operand
		default:
			panic(fmt.Sprintf("ERROR LessThan: unsupported type of field value for key %s", field))
		}
	}, fmt.Sprintf("%s < %.2f", field, operand))
}

func (r *RuleEngine[T]) LessThanEqualTo(field string, operand float64) *SimpleCondition[T] {
	return r.CustomFn(func(data T) bool {
		fieldValue := reflect.ValueOf(data).FieldByName(field)
		switch fieldValue.Kind() {
		case reflect.Int, reflect.Int32, reflect.Int64:
			return fieldValue.Int() <= int64(operand)
		case reflect.Float32, reflect.Float64:
			return fieldValue.Float() <= operand
		default:
			panic(fmt.Sprintf("ERROR LessThanEqualTo: unsupported type of field value for key %s", field))
		}
	}, fmt.Sprintf("%s <= %.2f", field, operand))
}

func (r *RuleEngine[T]) EqualToNum(field string, value float64) *SimpleCondition[T] {
	return r.CustomFn(func(data T) bool {
		fieldValue := reflect.ValueOf(data).FieldByName(field)

		switch fieldValue.Kind() {
		case reflect.Int, reflect.Int32, reflect.Int64:
			return fieldValue.Int() == int64(value)
		case reflect.Float32, reflect.Float64:
			return fieldValue.Float() == value
		default:
			panic(fmt.Sprintf("ERROR EqualToNum: unsupported type of field value for key %s", field))
		}
	}, fmt.Sprintf("%s == %.2f", field, value))
}

func (r *RuleEngine[T]) EqualToStr(field string, value string) *SimpleCondition[T] {
	return r.CustomFn(func(data T) bool {
		fieldValue := reflect.ValueOf(data).FieldByName(field)

		switch fieldValue.Kind() {
		case reflect.String:
			return fieldValue.String() == value
		default:
			panic(fmt.Sprintf("ERROR EqualToStr: unsupported type of field value for key %s", field))
		}
	}, fmt.Sprintf("%s == %s", field, value))
}

func (r *RuleEngine[T]) NotEqualToNum(field string, value float64) *SimpleCondition[T] {
	return r.CustomFn(func(data T) bool {
		fieldValue := reflect.ValueOf(data).FieldByName(field)

		switch fieldValue.Kind() {
		case reflect.Int, reflect.Int32, reflect.Int64:
			return fieldValue.Int() != int64(value)
		case reflect.Float32, reflect.Float64:
			return fieldValue.Float() != value
		default:
			panic(fmt.Sprintf("ERROR NotEqualToNum: unsupported type of field value for key %s", field))
		}
	}, fmt.Sprintf("%s != %.2f", field, value))
}

func (r *RuleEngine[T]) NotEqualToStr(field string, value string) *SimpleCondition[T] {
	return r.CustomFn(func(data T) bool {
		fieldValue := reflect.ValueOf(data).FieldByName(field)

		switch fieldValue.Kind() {
		case reflect.String:
			return fieldValue.String() != value
		default:
			panic(fmt.Sprintf("ERROR NotEqualToStr: unsupported type of field value for key %s", field))
		}
	}, fmt.Sprintf("%s != %s", field, value))
}
