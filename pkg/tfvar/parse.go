package tfvar

func ParseValues(from map[string]UnparsedVariableValue, vars []Variable) ([]Variable, error) {
	for i, v := range vars {
		unparsed, found := from[v.Name]
		if !found {
			continue
		}

		val, err := unparsed.ParseVariableValue(v.parsingMode)
		if err != nil {
			return nil, err
		}

		vars[i].Value = val
	}

	return vars, nil
}
