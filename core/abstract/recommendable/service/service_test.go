package service_test

import (
	"github.com/k0marov/go-socnet/core/abstract/recommendable/service"
	"github.com/k0marov/go-socnet/core/general/core_values"
	. "github.com/k0marov/go-socnet/core/helpers/test_helpers"
	"testing"
)

func TestRecsGetter(t *testing.T) {
	target := RandomId()
	count := 5
	recs := []string{RandomId(), RandomId(), RandomId(), RandomId(), RandomId()}
	t.Run("error case - storeRecsGetter throws", func(t *testing.T) {
		err := RandomError()
		storeRecsGetter := func(core_values.UserId, int) ([]string, error) {
			return nil, err
		}
		_, gotErr := service.NewRecsGetter(storeRecsGetter, nil)(target, count)
		AssertError(t, gotErr, err)
	})
	t.Run("happy case - storeRecsGetter returns enough recs", func(t *testing.T) {
		storeRecsGetter := func(targetId core_values.UserId, gotCount int) ([]string, error) {
			if targetId == target && gotCount == count {
				return recs, nil
			}
			panic("unexpected args")
		}
		got, err := service.NewRecsGetter(storeRecsGetter, nil)(target, count)
		AssertNoError(t, err)
		Assert(t, got, recs, "returned recommendations")
	})
	t.Run("storeRecsGetter returns not enough recs", func(t *testing.T) {
		returnedRecs := recs[:2]
		storeRecsGetter := func(core_values.UserId, int) ([]string, error) {
			return returnedRecs, nil
		}
		t.Run("happy case", func(t *testing.T) {
			randomRecs := recs[2:]
			randomRecsGetter := func(gotCount int) ([]string, error) {
				if gotCount == count-2 {
					return randomRecs, nil
				}
				panic("unexpected")
			}
			got, err := service.NewRecsGetter(storeRecsGetter, randomRecsGetter)(target, count)
			AssertNoError(t, err)
			Assert(t, got, recs, "merged recommendations")
		})
		t.Run("error case - getting random recs throws", func(t *testing.T) {
			randomRecsGetter := func(int) ([]string, error) {
				return nil, RandomError()
			}
			_, err := service.NewRecsGetter(storeRecsGetter, randomRecsGetter)(target, count)
			AssertSomeError(t, err)
		})
	})
}
