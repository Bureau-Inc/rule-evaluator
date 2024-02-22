package rule_evaluator

type Action struct {
	Fn          func(results map[string]interface{}) error
	Description string
}

type ActionExplanation struct {
	ActionDescription     string
	ConditionExplanations string
}

type Rule[T any] struct {
	Conditions Condition[T]
	Action     Action
}

type RuleEngine[T any] struct {
	rules        []Rule[T]
	Explanations []ActionExplanation
	results      map[string]interface{}
}

func (r *RuleEngine[T]) CreateSimpleCondition(condition bool, description string, invert bool) *SimpleCondition[T] {
	return &SimpleCondition[T]{
		Fn: func(data T) bool {
			if invert {
				return !condition
			}
			return condition
		},
		Description: description,
	}
}

func (r *RuleEngine[T]) CustomFn(Fn func(data T) bool, description string) *SimpleCondition[T] {
	return &SimpleCondition[T]{
		Fn:          Fn,
		Description: description,
	}
}

func (r *RuleEngine[T]) Not(condition bool, description string) *SimpleCondition[T] {
	return r.CreateSimpleCondition(condition, description, true)
}

func (r *RuleEngine[T]) Is(condition bool, description string) *SimpleCondition[T] {
	return r.CreateSimpleCondition(condition, description, false)
}

func (r *RuleEngine[T]) AnyOf(conditions ...Condition[T]) *ORCondition[T] {
	return &ORCondition[T]{
		Conditions: conditions,
	}
}

func (r *RuleEngine[T]) AllOf(conditions ...Condition[T]) *ANDCondition[T] {
	return &ANDCondition[T]{
		Conditions: conditions,
	}
}

func (c *RuleEngine[T]) DefineRule(condition Condition[T], actionFn func(results map[string]interface{}) error, actionDescription string) Rule[T] {
	return Rule[T]{
		Conditions: condition,
		Action: Action{
			Fn:          actionFn,
			Description: actionDescription,
		},
	}
}

func (r *RuleEngine[T]) AddRules(rules ...Rule[T]) {
	for _, rule := range rules {
		r.rules = append(r.rules, rule)
	}
}

type ResultsInitializerFunc func() map[string]interface{}

func (r *RuleEngine[T]) FireRules(data T, initializer ResultsInitializerFunc) (map[string]interface{}, error) {
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
