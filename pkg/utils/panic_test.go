package utils

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoErr(t *testing.T) {
	t.Parallel()

	t.Run("should not panic when error is nil", func(t *testing.T) {
		t.Parallel()
		assert.NotPanics(t, func() {
			NoErr(nil)
		})
	})

	t.Run("should panic when error is not nil", func(t *testing.T) {
		t.Parallel()
		errTest := errors.New("db connection failed")
		assert.PanicsWithValue(t, errTest, func() {
			NoErr(errTest)
		})
	})
}

func TestNoErrWithMsg(t *testing.T) {
	t.Parallel()

	t.Run("should not panic when error is nil", func(t *testing.T) {
		t.Parallel()
		assert.NotPanics(t, func() {
			NoErrWithMsg(nil, "custom error message")
		})
	})

	t.Run("should panic with formatted message when error is not nil", func(t *testing.T) {
		t.Parallel()
		errTest := errors.New("timeout")
		expectedPanicValue := "failed to load config: timeout"

		assert.PanicsWithValue(t, expectedPanicValue, func() {
			NoErrWithMsg(errTest, "failed to load config")
		})
	})
}
