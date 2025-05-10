package rules

import (
	"os"
	"strconv"

	"gopkg.in/yaml.v3"
)

func LoadRulesMap(path string) (map[string]Rule, int, error) {
	rulesMap := make(map[string]Rule)
	maxID := 0

	if _, err := os.Stat(path); err == nil {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, 0, err
		}

		var temp map[string]Rule
		if err := yaml.Unmarshal(data, &temp); err != nil {
			return nil, 0, err
		}

		for idStr, rule := range temp {
			// Обновляем максимальный ID
			id, err := strconv.Atoi(idStr)
			if err == nil && id > maxID {
				maxID = id
			}

			// Компилируем правило
			if err := rule.Compile(); err != nil {
				return nil, 0, err
			}
			rulesMap[idStr] = rule
		}
	}

	return rulesMap, maxID + 1, nil
}
