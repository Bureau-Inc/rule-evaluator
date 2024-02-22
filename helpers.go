package rule_evaluator

import (
	"fmt"
	"reflect"
)

// helper methods to create common conditions for a particular field in data in concise way

func (r *RuleEngine[T]) GreaterThan(field string, operand interface{}) *SimpleCondition[T] {
	return r.CustomFn(func(data T) bool {
		fieldValue := reflect.ValueOf(data).FieldByName(field)
		switch fieldValue.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return fieldValue.Int() > operand.(int64)
		case reflect.Float32, reflect.Float64:
			return fieldValue.Float() > operand.(float64)
		default:
			fmt.Printf("ERROR GreaterThan: unsupported type of field value for key %s", field)
			return false
		}
	}, fmt.Sprintf("%s > %v", field, operand))
}

func (r *RuleEngine[T]) GreaterThanEqualTo(field string, operand interface{}) *SimpleCondition[T] {
	return r.CustomFn(func(data T) bool {
		fieldValue := reflect.ValueOf(data).FieldByName(field)
		switch fieldValue.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return fieldValue.Int() >= operand.(int64)
		case reflect.Float32, reflect.Float64:
			return fieldValue.Float() >= operand.(float64)
		default:
			fmt.Printf("ERROR GreaterThanEqualTo: unsupported type of field value for key %s", field)
			return false
		}
	}, fmt.Sprintf("%s > %v", field, operand))
}

func (r *RuleEngine[T]) LessThan(field string, operand interface{}) *SimpleCondition[T] {
	return r.CustomFn(func(data T) bool {
		fieldValue := reflect.ValueOf(data).FieldByName(field)
		switch fieldValue.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return fieldValue.Int() < operand.(int64)
		case reflect.Float32, reflect.Float64:
			return fieldValue.Float() < operand.(float64)
		default:
			fmt.Printf("ERROR LessThan: unsupported type of field value for key %s", field)
			return false
		}
	}, fmt.Sprintf("%s > %v", field, operand))
}

func (r *RuleEngine[T]) LessThanEqualTo(field string, operand interface{}) *SimpleCondition[T] {
	return r.CustomFn(func(data T) bool {
		fieldValue := reflect.ValueOf(data).FieldByName(field)
		switch fieldValue.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return fieldValue.Int() <= operand.(int64)
		case reflect.Float32, reflect.Float64:
			return fieldValue.Float() <= operand.(float64)
		default:
			fmt.Printf("ERROR LessThanEqualTo: unsupported type of field value for key %s", field)
			return false
		}
	}, fmt.Sprintf("%s > %v", field, operand))
}

func (r *RuleEngine[T]) EqualTo(field string, value interface{}) *SimpleCondition[T] {
	return r.CustomFn(func(data T) bool {
		fieldValue := reflect.ValueOf(data).FieldByName(field)

		switch fieldValue.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return fieldValue.Int() == value.(int64)
		case reflect.Float32, reflect.Float64:
			return fieldValue.Float() == value.(float64)
		case reflect.String:
			return fieldValue.String() == value.(string)
		default:
			fmt.Printf("ERROR EqualTo: unsupported type of field value for key %s", field)
			return false
		}
	}, fmt.Sprintf("%s == %v", field, value))
}

func (r *RuleEngine[T]) NotEqualTo(field string, value interface{}) *SimpleCondition[T] {
	return r.CustomFn(func(data T) bool {
		fieldValue := reflect.ValueOf(data).FieldByName(field)

		switch fieldValue.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return fieldValue.Int() != value.(int64)
		case reflect.Float32, reflect.Float64:
			return fieldValue.Float() != value.(float64)
		case reflect.String:
			return fieldValue.String() != value.(string)
		default:
			fmt.Printf("ERROR NotEqualTo: unsupported type of field value for key %s", field)
			return false
		}
	}, fmt.Sprintf("%s == %v", field, value))
}
