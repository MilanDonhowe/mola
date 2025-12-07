package syntax

import (
	"fmt"
	"mola/language"
	"strconv"
	"strings"
)

func Pr_Str(mal language.MalValue) (string, error) {
	switch mal.TypeId {
	case language.Symbol:
		return mal.Symbol, nil
	case language.Integer:
		return strconv.Itoa(mal.Integer), nil
	case language.List:
		var p []string
		for _, v := range mal.List {
			s, err := Pr_Str(v)
			if err != nil {
				return "", err
			}
			p = append(p, s)
		}
		return "(" + strings.Join(p, " ") + ")", nil
	case language.Nil:
		return "nil", nil
	case language.String:
		return mal.Symbol, nil
	}
	return "", fmt.Errorf("no string representation logic for supplised mal type id=%d", mal.TypeId)
}
