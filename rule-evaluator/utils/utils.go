package utils

import (
	"regexp"
	"strings"
)

func FormatExplanation(explanation string) string {
	formattedExplanation := ""
	indent := 0

	for _, char := range explanation {
		switch char {
		case '[':
			formattedExplanation += "\n" + strings.Repeat("  ", indent) + "["
			indent++
		case ']':
			indent--
			formattedExplanation += "\n" + strings.Repeat("  ", indent) + "]"
		default:
			formattedExplanation += string(char)
		}
	}

	return formattedExplanation
}

// The following function is AI generated, because I don't know regex :)
func ParseRuleString(ruleStr string) ([][][3]string, error) {
	// Regex to extract individual conditions
	conditionRegex := regexp.MustCompile(`([a-zA-Z0-9_]+)\s*([><=!]+)\s*([^;]+)`)

	// Remove outer brackets and split on 'AND:'
	innerStr := strings.TrimSpace(ruleStr[1 : len(ruleStr)-1])
	andConditions := strings.Split(innerStr, "AND:")

	result := make([][][3]string, 0)
	for _, andClause := range andConditions {
		conditions := conditionRegex.FindAllStringSubmatch(andClause, -1)
		if conditions == nil {
			continue
		}

		clauseResult := make([][3]string, 0)
		for _, cond := range conditions {
			clauseResult = append(clauseResult, [3]string{cond[1], cond[2], cond[3]})
		}
		result = append(result, clauseResult)
	}

	return result, nil
}
