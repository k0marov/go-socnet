package validators_test

import (
	"github.com/k0marov/go-socnet/core/general/client_errors"
	. "github.com/k0marov/go-socnet/core/helpers/test_helpers"
	"strings"
	"testing"

	"github.com/k0marov/go-socnet/features/comments/domain/validators"
	"github.com/k0marov/go-socnet/features/comments/domain/values"
)

func TestCommentValidator(t *testing.T) {
	cases := []struct {
		comment values.NewCommentValue

		isValid bool
		wantErr client_errors.ClientError
	}{
		{values.NewCommentValue{Text: "Normal text"}, true, client_errors.ClientError{}},
		{values.NewCommentValue{Text: ""}, false, client_errors.EmptyText},
		{values.NewCommentValue{Text: strings.Repeat("looooong", 100)}, false, client_errors.TextTooLong},
	}

	for _, testCase := range cases {
		t.Run(testCase.comment.Text, func(t *testing.T) {
			gotErr, gotValid := validators.NewCommentValidator()(testCase.comment)
			AssertFatal(t, gotValid, testCase.isValid, "the result of validation")
			Assert(t, gotErr, testCase.wantErr, "the returned client error")
		})
	}
}
