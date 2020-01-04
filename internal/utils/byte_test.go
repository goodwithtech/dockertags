package utils

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestByteSize(t *testing.T) {
	testcases := map[string]struct {
		size     int
		expected string
	}{
		"AllowMinusNum": {
			size:     -1024,
			expected: "-1024",
		},
		"0": {
			size:     0,
			expected: "0B",
		},
		"Byte": {
			size:     10,
			expected: "10B",
		},
		"1.5K": {
			size:     1024 + 512,
			expected: "1.5K",
		},
		"1.7499..K": {
			size:     1024 + 512 + 255,
			expected: "1.7K",
		},
		"1.75K": {
			size:     1024 + 512 + 256,
			expected: "1.8K",
		},
		"1MB": {
			size:     1024 * 1024,
			expected: "1M",
		},
		"1GB": {
			size:     1024 * 1024 * 1024,
			expected: "1G",
		},
		"1TB": {
			size:     1024 * 1024 * 1024 * 1024,
			expected: "1024G",
		},
	}
	for tc, v := range testcases {
		actual := ByteSize(v.size)
		if diff := cmp.Diff(actual, v.expected); diff != "" {
			t.Errorf("%s: diff %v", tc, diff)
		}
	}
}
