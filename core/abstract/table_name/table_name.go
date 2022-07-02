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

func NewTableName(name string) TableName {
	if name != "" && len(strings.FieldsFunc(name, unicode.IsLetter)) == 0 {
		return TableName{name, true}
	} else {
		return TableName{"", false}
	}
}

func (t TableName) Value() (string, error) {
	if t.isValid {
		return t.value, nil
	} else {
		return "", fmt.Errorf("the provided table name is not valid")
	}
}
