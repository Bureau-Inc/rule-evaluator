package rule_evaluator

import (
	"fmt"
	"strings"
)

type Condition interface {
	Evaluate() bool
	GetDescription() string
}

type SimpleCondition struct {
	Fn          func() bool
	Description string
}

type CompositeCondition struct {
	Type       string
	Conditions []Condition
}

func (s SimpleCondition) Evaluate() bool {
	return s.Fn()
}

func (s SimpleCondition) GetDescription() string {
	return s.Description
}

func (c CompositeCondition) Evaluate() bool {
	if c.Type == "OR" {
		for _, cond := range c.Conditions {
			result := cond.Evaluate()
			if result {
				return true
			}
		}
		return false
	} else if c.Type == "AND" {
		var descriptions []string

		for _, cond := range c.Conditions {
			result := cond.Evaluate()

			if !result {
				return false
			}
			descriptions = append(descriptions, c.GetDescription())
		}

		return true
	}
	fmt.Println("Error: This condition type is not supported: ", c.Type)
	return false
}

func (c CompositeCondition) GetDescription() string {
	// Example: Combine descriptions of nested conditions if needed
	return fmt.Sprintf("[%s: %s]", c.Type, combineDescriptions(getDescriptions(c)))
}

// Helper for getting descriptions from a nested condition
func getDescriptions(cond Condition) []string {
	if comp, ok := cond.(CompositeCondition); ok {
		var descriptions []string
		for _, nestedCond := range comp.Conditions {
			descriptions = append(descriptions, nestedCond.GetDescription())
		}
		return descriptions
	} else if simpl, ok := cond.(SimpleCondition); ok {
		// If it's a base condition (likely a SimpleCondition)
		return []string{simpl.GetDescription()}
	}
	fmt.Println("Error in description: This condition type is not supported ")
	return []string{}
}

func combineDescriptions(descriptions []string) string {
	return strings.Join(descriptions, "; ")
}

type ActionFunc func(results *Results) (string, error)

type ActionExplanation struct {
	ActionDescription     string
	ConditionExplanations string
}

type Results struct {
	Data         map[string]interface{}
	Explanations []ActionExplanation
}

type Rule struct {
	Conditions Condition
	Action     ActionFunc
}

func CreateSimpleCondition(condition bool, description string, invert bool) SimpleCondition {
	return SimpleCondition{
		Fn: func() bool {
			if invert {
				return !condition
			}
			return condition
		},
		Description: description,
	}
}

func CustomFn(Fn func() bool, description string) SimpleCondition {
	return SimpleCondition{
		Fn:          Fn,
		Description: description,
	}
}

func Not(condition bool, description string) SimpleCondition {
	return CreateSimpleCondition(condition, description, true)
}

func Is(condition bool, description string) SimpleCondition {
	return CreateSimpleCondition(condition, description, false)
}

func AnyOf(conditions ...Condition) CompositeCondition {
	return CompositeCondition{
		Type:       "OR",
		Conditions: conditions,
	}
}

func AllOf(conditions ...Condition) CompositeCondition {
	return CompositeCondition{
		Type:       "AND",
		Conditions: conditions,
	}
}

type ResultsInitializerFunc func() map[string]interface{}

func FireRules(rules []Rule, initializer ResultsInitializerFunc) (*Results, error) {
	results := &Results{Data: initializer(), Explanations: make([]ActionExplanation, 0)}

	for _, ruleItem := range rules {
		ruleResult := ruleItem.Conditions.Evaluate()

		if ruleResult {
			desc, err := ruleItem.Action(results)
			if err != nil {
				return nil, err
			}
			results.Explanations = append(results.Explanations, ActionExplanation{
				ActionDescription:     desc,
				ConditionExplanations: ruleItem.Conditions.GetDescription(),
			})
		}
	}

	return results, nil
}
