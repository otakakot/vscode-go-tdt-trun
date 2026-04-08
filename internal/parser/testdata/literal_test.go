package testdata

import "testing"

func TestLiteral(t *testing.T) {
	t.Run("first subtest", func(t *testing.T) {
		_ = 1
	})
	t.Run("second subtest", func(t *testing.T) {
		_ = 2
	})
}
