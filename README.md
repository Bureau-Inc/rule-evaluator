# rule-evaluator

This project provides a rule engine for evaluating dynamic rules. It can be used to uncover the specific conditions that drive actions in complex rule-based systems.

## Features

Define rules using conditions like:

- Get description of conditions that drove a particular action
- Built in conditions like GreaterThan (>), LessThan (<), EqualTo (==), In (for slices, Nil checks
- Combine conditions using AND and OR
- Write custom functions for rules if you have more complex conditions
- Specify actions to be executed when rules are met
- All conditions and actions are just usual Go functions

## Example Usage

Say you would like to label a user session as risky based on some conditions that you have stored in some configuration. And the conditions look like this:

IF:

- user.money > 2000 _AND_ user.employment = nil

_OR_

- one of the user's friends is "John"

Then:

- Mark the session category as risky

Here is how you would go about doing it (see complete code in `examples/risk_rules.go`):

```

import (
	"fmt"
	r "github.com/Bureau-Inc/rule-evaluator/rule-evaluator"
)


func AnalyseUserRisk() error {

	userData := UserData{
		Name:       "Paul",
		Age:        29,
		Money:      500000,
		Employment: nil,
		Friends:    []string{"Saul"},
	}

    // Step 1: Create a RuleEngine.
	re := r.RuleEngine{}

    // Step2: Define the above mentioned rules and the corresponding action
	highRiskRule := re.DefineRule(
        // conditions
		re.OR(
			re.AND(re.GreaterThan("Money", 20000), re.IsNil("Employment")),
			re.ListContainsStr("Friends", "John"),
		),

        // action function (UpdateRiskScore returns a function)
		UpdateRiskScore("highRiskCategory"),

        // action description
		fmt.Sprintf("risk level: highRiskCategory"),
	)

	// add as many rules as you want
	re.AddRules(highRiskRule)

    // Step 3: fire the rules by passing the input data and results initializer function
	results, err := re.FireRules(userData, InitializeSessionResultData)

	fmt.Printf("risk category: %s \n", results["riskScore"])


    // Step 4: Inspect the results
	fmt.Printf("%s", re.InspectLastSession())
	return nil

}

```

Output:

```
[AND: Money > 20000.00; Employment == nil
]
 Action-> risk level: highRiskCategory
```

## Getting Started

Install the package: `go get github.com/Bureau-Inc/rule-evaluator/rule-evaluator`

Refer to the example for basic usage.
Extend the engine with custom conditions and actions.
