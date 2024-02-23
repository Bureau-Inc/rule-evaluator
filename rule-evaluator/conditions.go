package rule_evaluator

import (
	"fmt"
	"strings"
)

type Condition[T any] interface {
	Evaluate(data T) bool
	GetDescription() string
}

type SimpleCondition[T any] struct {
	Fn          func(data T) bool
	Description string
}

type ORCondition[T any] struct {
	Conditions       []Condition[T]
	WinningCondition Condition[T] // for generating runtime explanation
}

type ANDCondition[T any] struct {
	Conditions []Condition[T]
}

func (s *SimpleCondition[T]) Evaluate(data T) bool {
	return s.Fn(data)
}

func (s *SimpleCondition[T]) GetDescription() string {
	return s.Description
}

func (c *ORCondition[T]) Evaluate(data T) bool {
	for _, cond := range c.Conditions {
		result := cond.Evaluate(data)
		if result {
			c.WinningCondition = cond
			return true
		}
	}
	return false
}

func (c *ORCondition[T]) GetDescription() string {
	return c.WinningCondition.GetDescription()
}

func (c *ORCondition[T]) GetStaticDescription() string {

	var descriptions []string
	for _, nestedCond := range c.Conditions {
		descriptions = append(descriptions, nestedCond.GetDescription())
	}
	return strings.Join(descriptions, "; ")
}

func (c *ANDCondition[T]) Evaluate(data T) bool {
	for _, cond := range c.Conditions {
		result := cond.Evaluate(data)
		if !result {
			return false
		}
	}
	return true
}

func (c *ANDCondition[T]) GetDescription() string {
	var descriptions []string
	for _, nestedCond := range c.Conditions {
		descriptions = append(descriptions, nestedCond.GetDescription())
	}
	return fmt.Sprintf("[AND: %s]", strings.Join(descriptions, "; "))
}
