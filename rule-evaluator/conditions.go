package rule_evaluator

import (
	"fmt"
	"strings"
)

type Condition interface {
	Evaluate(data interface{}) bool
	GetDescription() string
}

type SimpleCondition struct {
	Fn          func(data interface{}) bool
	Description string
}

type ORCondition struct {
	Conditions       []Condition
	WinningCondition Condition // for generating runtime explanation
}

type ANDCondition struct {
	Conditions []Condition
}

func (s *SimpleCondition) Evaluate(data interface{}) bool {
	return s.Fn(data)
}

func (s *SimpleCondition) GetDescription() string {
	return s.Description
}

func (c *ORCondition) Evaluate(data interface{}) bool {
	for _, cond := range c.Conditions {
		result := cond.Evaluate(data)
		if result {
			c.WinningCondition = cond
			return true
		}
	}
	return false
}

func (c *ORCondition) GetDescription() string {
	return c.WinningCondition.GetDescription()
}

func (c *ORCondition) GetStaticDescription() string {

	var descriptions []string
	for _, nestedCond := range c.Conditions {
		descriptions = append(descriptions, nestedCond.GetDescription())
	}
	return strings.Join(descriptions, "; ")
}

func (c *ANDCondition) Evaluate(data interface{}) bool {
	for _, cond := range c.Conditions {
		result := cond.Evaluate(data)
		if !result {
			return false
		}
	}
	return true
}

func (c *ANDCondition) GetDescription() string {
	var descriptions []string
	for _, nestedCond := range c.Conditions {
		descriptions = append(descriptions, nestedCond.GetDescription())
	}
	return fmt.Sprintf("[AND: %s]", strings.Join(descriptions, "; "))
}
