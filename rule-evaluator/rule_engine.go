package rule_evaluator

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

func (r *RuleEngine) DefineRule(condition Condition, action Action) Rule {
	return Rule{
		Conditions: condition,
		Action:     action,
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
