package rule_evaluator_test

import (
	"fmt"
	r "github.com/Bureau-Inc/rule-evaluator/rule-evaluator"
	"github.com/Bureau-Inc/rule-evaluator/rule-evaluator/utils"
	"testing"
)

// Sample Data Structure - Make sure this matches yours
type UserData struct {
	Name       string      `json:"name"`
	Age        int         `json:"age"`
	Money      float64     `json:"money"`
	Employment interface{} `json:"employment"`
	Friends    []string    `json:"friends"`
}

func InitializeSessionResultData() map[string]interface{} {
	// Your initialization logic here
	return map[string]interface{}{
		"riskScore": "lowRiskCategory",
	}
}

func UpdateRiskScore(riskCat string) r.ActionFn {
	return func(results map[string]interface{}) error {
		results["riskScore"] = riskCat
		return nil
	}
}

func areAllThreeEqual(a, b [3]string) bool {
	return a[0] == b[0] && a[1] == b[1] && a[2] == b[2]
}

func TestHighRiskRule(t *testing.T) {
	userData := UserData{
		Name:       "Jesse",
		Age:        29,
		Money:      500000,
		Employment: nil,
		Friends:    []string{"Saul", "Mike", "Walter"},
	}

	re := r.RuleEngine{}
	highRiskRule := re.DefineRule(
		re.OR(
			re.AND(re.GreaterThan("Money", 20000), re.IsNil("Employment")),
			re.ListContainsStr("Friends", "Tuco"),
		),
		UpdateRiskScore("highRiskCategory"),
		fmt.Sprintf("risk level: highRiskCategory"),
	)
	re.AddRules(highRiskRule)

	results, err := re.FireRules(userData, InitializeSessionResultData)
	if err != nil {
		t.Fatalf("Error evaluating rules: %v", err)
	}

	if results["riskScore"] != "highRiskCategory" {
		t.Errorf("Expected riskScore 'highRiskCategory', got: %v", results["riskScore"])
	}

	descriptionStr := re.Explanations[0].ConditionExplanations
	parsed_expl, err := utils.ParseRuleString(descriptionStr)
	if err != nil {
		t.Fatalf("Error in parseRuleString, %v", err)
	}
	moneyCondOk := areAllThreeEqual(parsed_expl[0][0], [3]string{"Money", ">", "20000.00"})
	if !moneyCondOk {
		t.Errorf("Condition description not okay: %v", parsed_expl[0][0])
	}
	employCondOk := areAllThreeEqual(parsed_expl[0][1], [3]string{"Employment", "==", "nil"})
	if !employCondOk {
		t.Errorf("Condition description not okay: %v", parsed_expl[0][1])
	}

}
