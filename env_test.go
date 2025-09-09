package common_test

import (
	"os"
	"testing"

	common "github.com/SoeltanIT/agg-common-be"
	"github.com/stretchr/testify/assert"
)

func TestGetEnv(t *testing.T) {
	t.Run("GetEnv returns a value from the environment", func(t *testing.T) {
		_ = os.Setenv("TEST_ENV", "test")
		result := common.GetEnv("TEST_ENV", "fallback")
		assert.Equal(t, "test", result, "GetEnv() should return the value from the environment")
		_ = os.Unsetenv("TEST_ENV")
	})

	t.Run("GetEnv returns a fallback value if the environment variable is not set", func(t *testing.T) {
		result := common.GetEnv("TEST_ENV", "fallback")
		assert.Equal(t, "fallback", result, "GetEnv() should return the fallback value if the environment variable is not set")
	})
}
