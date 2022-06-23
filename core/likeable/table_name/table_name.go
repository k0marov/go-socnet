package table_name

import (
	"fmt"
	"strings"
	"unicode"
)

type TableName struct {
	value   string
	isValid bool
}

func NewTableName(name string) (TableName, error) {
	if name != "" && len(strings.FieldsFunc(name, unicode.IsLetter)) == 0 {
		return TableName{name, true}, nil
	} else {
		return TableName{"", false}, fmt.Errorf("the provided table name is not valid: %s", name)
	}
}

func (t TableName) Value() (string, error) {
	if t.isValid {
		return t.value, nil
	} else {
		return "", fmt.Errorf("the provided table name is not valid")
	}
}
