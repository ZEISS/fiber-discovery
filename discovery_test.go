package discovery_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	discovery "github.com/zeiss/fiber-discovery"
)

func TestWellknownURL(t *testing.T) {
	tests := []struct {
		keys     []string
		expected string
	}{
		{
			keys:     []string{"test"},
			expected: ".well-known/test",
		},
		{
			keys:     []string{"test", "test2"},
			expected: ".well-known/test/test2",
		},
		{
			keys:     []string{"test", "test2", "test3"},
			expected: ".well-known/test/test2/test3",
		},
	}

	for _, test := range tests {
		result := discovery.WellknownURL(test.keys...)
		require.Equal(t, test.expected, result)
	}
}
