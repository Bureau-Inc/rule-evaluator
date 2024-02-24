package examples

import (
	"fmt"
	r "github.com/Bureau-Inc/rule-evaluator/rule-evaluator"
)

func InitializeSessionResultData() map[string]interface{} {
	return map[string]interface{}{
		"riskScore":    0.0,
		"customerType": "Normal",
	}
}

func UpdateRiskScore(riskCat string) r.ActionFn {
	return func(results map[string]interface{}) error {
		results["riskScore"] = riskCat
		return nil
	}
}

func UpdateCustomerType(customerType string) r.ActionFn {
	return func(results map[string]interface{}) error {
		results["customerType"] = customerType
		return nil
	}
}

type EmploymentDetail struct {
	Company string
	Salary  string
}

type UserData struct {
	Name       string
	Age        int
	Money      float64
	Employment *EmploymentDetail
	Friends    []string
	OrderCount int
	PlanType   string
}

func AnalyseUserRisk() error {

	userData := UserData{
		Name:       "Jesse",
		Age:        29,
		Money:      500000,
		Employment: nil,
		Friends:    []string{"Saul", "Mike", " Walter"},
		OrderCount: 100,
		PlanType:   "Basic",
	}

	re := r.RuleEngine{}
	highValueCustomerRule := re.DefineRule(
		re.CustomFn(func(data interface{}) bool {
			userData, ok := data.(UserData)
			if !ok {
				return false
			}
			// Logic to count orders for UserID within the time window ...
			return userData.OrderCount >= 50
		}, "High order count"),
		UpdateCustomerType("HighValue"),
		"customer type = HighValue")

	highRiskRule := re.DefineRule(
		re.OR(
			re.AND(re.GreaterThan("Money", 20000), re.IsNil("Employment")),
			re.ListContainsStr("Friends", "Tuco"),
		),
		UpdateRiskScore("highRiskCategory"),
		fmt.Sprintf("risk level: highRiskCategory"),
	)

	re.AddRules(highValueCustomerRule, highRiskRule)

	results, err := re.FireRules(userData, InitializeSessionResultData)
	if err != nil {
		return fmt.Errorf("error in evaluating rules")
	}
	fmt.Printf("risk category: %s \n", results["riskScore"])
	fmt.Printf("%s", re.InspectLastSession())
	return nil
}
