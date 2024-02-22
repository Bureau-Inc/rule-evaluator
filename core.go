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

type ORCondition struct {
	Conditions       []Condition
	WinningCondition Condition // for generating runtime explanation
}

type ANDCondition struct {
	Conditions      []Condition
	LosingCondition Condition
}

func (s *SimpleCondition) Evaluate() bool {
	return s.Fn()
}

func (s *SimpleCondition) GetDescription() string {
	return s.Description
}

func (c *ORCondition) Evaluate() bool {
	for _, cond := range c.Conditions {
		result := cond.Evaluate()
		if result {
			c.WinningCondition = cond
			return true
		}
	}
	return false

	//else if c.Type == "AND" {
	//	var descriptions []string
	//
	//	for _, cond := range c.Conditions {
	//		result := cond.Evaluate()
	//		if !result {
	//			c.Explanation = ""
	//			return false
	//		}
	//		descriptions = append(descriptions, GetExplanationForCondition(cond))
	//	}
	//	c.Explanation = combineDescriptions(descriptions)
	//	return true
	//}

}

func (c *ANDCondition) Evaluate() bool {
	for _, cond := range c.Conditions {
		result := cond.Evaluate()
		if !result {
			return false
		}
	}
	return true
}

// [AND: User is Subscribed; [AND: User is PM; User is rich]] -> action

func (c *ANDCondition) GetDescription() string {
	// Example: Combine descriptions of nested conditions if needed
	return fmt.Sprintf("[AND: %s]", getDescription(c))
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

func getDescription(cond Condition) string {
	if and, ok := cond.(*ANDCondition); ok {
		var descriptions []string
		for _, nestedCond := range and.Conditions {
			descriptions = append(descriptions, nestedCond.GetDescription())
		}
		return strings.Join(descriptions, "; ")
	} else if or, ok := cond.(*ORCondition); ok {
		return or.GetDescription()
	} else if simpl, ok := cond.(*SimpleCondition); ok {
		return simpl.GetDescription()
	}
	fmt.Println("Error in description: This condition type is not supported ")
	return ""
}

type Action struct {
	Fn          func(results *Results) error
	Description string
}

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
	Action     Action
}

func CreateSimpleCondition(condition bool, description string, invert bool) *SimpleCondition {
	return &SimpleCondition{
		Fn: func() bool {
			if invert {
				return !condition
			}
			return condition
		},
		Description: description,
	}
}

func CustomFn(Fn func() bool, description string) *SimpleCondition {
	return &SimpleCondition{
		Fn:          Fn,
		Description: description,
	}
}

func Not(condition bool, description string) *SimpleCondition {
	return CreateSimpleCondition(condition, description, true)
}

func Is(condition bool, description string) *SimpleCondition {
	return CreateSimpleCondition(condition, description, false)
}

func AnyOf(conditions ...Condition) *ORCondition {
	return &ORCondition{
		Conditions: conditions,
	}
}

func AllOf(conditions ...Condition) *ANDCondition {
	return &ANDCondition{
		Conditions: conditions,
	}
}

type ResultsInitializerFunc func() map[string]interface{}

func FireRules(rules []Rule, initializer ResultsInitializerFunc) (*Results, error) {
	results := &Results{Data: initializer(), Explanations: make([]ActionExplanation, 0)}

	for _, ruleItem := range rules {
		ruleResult := ruleItem.Conditions.Evaluate()

		if ruleResult {
			err := ruleItem.Action.Fn(results)
			if err != nil {
				return nil, err
			}
			results.Explanations = append(results.Explanations, ActionExplanation{
				ActionDescription:     ruleItem.Action.Description,
				ConditionExplanations: ruleItem.Conditions.GetDescription(),
			})
		}
	}

	return results, nil
}
