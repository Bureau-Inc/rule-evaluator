package rule_evaluator

import (
	"fmt"
	"github.com/Bureau-Inc/rule-evaluator/rule-evaluator/utils"
	"strings"
)

type ActionFn func(results map[string]interface{}) error
type Action struct {
	Fn          func(results map[string]interface{}) error
	Description string
}

type ActionExplanation struct {
	ActionDescription     string
	ConditionExplanations string
}

type Rule struct {
	Conditions Condition
	Action     Action
}

type RuleEngine struct {
	rules        []Rule
	Explanations []ActionExplanation
	results      map[string]interface{}
}

type RuleEngineError struct {
	ErrorType string
	Field     string
	Message   string
}

func (r *RuleEngine) CreateSimpleCondition(condition bool, description string, invert bool) *SimpleCondition {
	return &SimpleCondition{
		Fn: func(data interface{}) bool {
			if invert {
				return !condition
			}
			return condition
		},
		Description: description,
	}
}

func (c *RuleEngine) DefineRule(condition Condition, actionFn func(results map[string]interface{}) error, actionDescription string) Rule {
	return Rule{
		Conditions: condition,
		Action: Action{
			Fn:          actionFn,
			Description: actionDescription,
		},
	}
}

func (r *RuleEngine) AddRules(rules ...Rule) {
	for _, rule := range rules {
		r.rules = append(r.rules, rule)
	}
}

type ResultsInitializerFunc func() map[string]interface{}

func (r *RuleEngine) FireRules(data interface{}, initializer ResultsInitializerFunc) (map[string]interface{}, error) {
	r.results = initializer()
	r.Explanations = make([]ActionExplanation, 0)

	for _, ruleItem := range r.rules {
		ruleResult := ruleItem.Conditions.Evaluate(data)

		if ruleResult {
			err := ruleItem.Action.Fn(r.results)
			if err != nil {
				return nil, err
			}
			r.Explanations = append(r.Explanations, ActionExplanation{
				ActionDescription:     ruleItem.Action.Description,
				ConditionExplanations: ruleItem.Conditions.GetDescription(),
			})
		}
	}
	return r.results, nil
}

func (r *RuleEngine) InspectLastSession() string {
	var descriptions []string
	for _, expl := range r.Explanations {
		descriptions = append(descriptions,
			fmt.Sprintf("%s \n Action-> %s\n", utils.FormatExplanation(expl.ConditionExplanations), expl.ActionDescription))
	}
	return strings.Join(descriptions, "\n")
}
