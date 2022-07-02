package table_name_test

import (
	"github.com/k0marov/go-socnet/core/abstract/likeable/table_name"
	. "github.com/k0marov/go-socnet/core/helpers/test_helpers"
	"testing"
)

func TestTableName(t *testing.T) {
	cases := []struct {
		tableName string
		isValid   bool
	}{
		{"Profile", true},
		{"'; DROP TABLE Profile; --", false},
		{"", false},
	}

	for _, testCase := range cases {
		t.Run(testCase.tableName, func(t *testing.T) {
			tblName := table_name.NewTableName(testCase.tableName)
			tblNameValue, err := tblName.Value()
			if testCase.isValid {
				AssertNoError(t, err)
				Assert(t, tblNameValue, testCase.tableName, "stored table name")
			} else {
				AssertSomeError(t, err)
			}
		})
	}
}
