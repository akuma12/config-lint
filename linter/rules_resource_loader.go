package linter

import (
	"github.com/stelligent/config-lint/assertion"
)

// RulesResourceLoader converts a YAML configuration file into a collection with Resource objects
type RulesResourceLoader struct{}

func getAttr(m map[string]interface{}, keys ...string) []interface{} {
	for _, key := range keys {
		if r, ok := m[key].([]interface{}); ok {
			return r
		}
	}
	return []interface{}{}
}

// Load converts a text file into a collection of Resource objects
func (l RulesResourceLoader) Load(filename string) ([]assertion.Resource, error) {
	resources := make([]assertion.Resource, 0)
	yamlResources, err := loadYAML(filename)
	if err != nil {
		return resources, err
	}
	for _, ruleSet := range yamlResources {
		setResource := assertion.Resource{
			ID:         getResourceIDFromFilename(filename),
			Type:       "LintRuleSet",
			Properties: ruleSet,
			Filename:   filename,
		}
		resources = append(resources, setResource)
		// The LintRuleSet resources already has an attribute called Rules
		// but also adding each of these rules as a separate LintRule resource
		// makes writing rules a lot simpler
		m := ruleSet.(map[string]interface{})
		rules := getAttr(m, "rules", "Rules")
		for _, rule := range rules {
			properties := rule.(map[string]interface{})
			ruleResource := assertion.Resource{
				ID:         properties["id"].(string),
				Type:       "LintRule",
				Properties: properties,
				Filename:   filename,
			}
			resources = append(resources, ruleResource)
		}
	}
	return resources, nil
}